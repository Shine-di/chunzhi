package wrapper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitee.com/risewinter/data-lol/library/log"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"go.elastic.co/apm"
	"go.uber.org/zap"
)

func ApmWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		md, ok := metadata.FromContext(ctx)
		if !ok {
			md = make(map[string]string)
		}
		md = metadata.Copy(md)
		transaction := apm.DefaultTracer.StartTransaction(req.Method(), "grpc")
		transaction.Context.SetTag("traceId", fmt.Sprintf("%s", md["x-b3-traceid"]))
		transaction.Context.SetTag("args", fmt.Sprintf("%v", req.Body()))
		if transaction != nil {
			ctx = apm.ContextWithTransaction(ctx, transaction)
			defer transaction.End()
		}
		return fn(ctx, req, rsp)
	}
}
func ErrWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
		startTime := time.Now()
		defer func() {
			var errMsg string
			if errRecover := recover(); errRecover != nil {
				isError, ok := errRecover.(error)
				if ok {
					errMsg = isError.Error()
					err = isError
				} else {
					errMsg = fmt.Sprintf("%v", errRecover)
					err = errors.New(fmt.Sprintf("%v", errRecover))
				}

				md, ok := metadata.FromContext(ctx)
				if !ok {
					md = make(map[string]string)
				}

				md = metadata.Copy(md)
				traceId := md["x-b3-traceid"]
				e := apm.DefaultTracer.NewError(err)
				e.Context.SetTag("args", fmt.Sprintf("%v", req.Body()))
				e.Context.SetTag("traceId", traceId)
				e.Context.SetTag("errMsg", errMsg)
				e.Send()
			}
			log.InfoWithCtx(ctx, "end request", zap.Any("request", req.Body()), zap.String("method", req.Method()), zap.String("request_time", time.Since(startTime).String()), zap.Error(err), zap.String("err_msg", errMsg))
		}()
		log.InfoWithCtx(ctx, "receive request Record", zap.String("method", req.Method()), zap.Any("req", req.Body()))

		err = fn(ctx, req, rsp)
		return err
	}
}
