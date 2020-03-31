package grpcbase

import (
	"context"
	"errors"
	"game-test/library/config"
	"game-test/library/log"
	site_var "game-test/library/site-var"
	"game-test/library/zipkin"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	openlog "github.com/opentracing/opentracing-go/log"
	"github.com/processout/grpc-go-pool"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceMap    = make(map[string]*pollService, 0)
	lock          = &sync.Mutex{}
	lockConfigMap = &sync.Mutex{}
	INIT_NUM      = 1
	INIT_CAPACITY = 100
)

type ServiceConfig struct {
	serviceName string
	address     string
	init        int
	capacity    int
}

type pollService struct {
	pool *grpcpool.Pool
}

func WithClientInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(clientInterceptorWithZipkin())
}

func WithStreamInterceptor() grpc.DialOption {
	return grpc.WithStreamInterceptor(streamClientInterceptorWithZipkin())
}

func clientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		log.Error("Invoked Error", zap.String("method", method), zap.Error(err))
	}

	log.Info("Invoked RPC Record", zap.String("method", method), zap.String("Duration", time.Since(start).String()), zap.Error(err), zap.Any("req", req), zap.Any("reply", reply))
	return err
}

func clientInterceptorWithZipkin(optFuncs ...zipkin.Option) grpc.UnaryClientInterceptor {
	otgrpcOpts := zipkin.NewOptions()
	otgrpcOpts.Apply(optFuncs...)
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var err error
		var parentCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}

		if otgrpcOpts.InclusionFunc != nil &&
			!otgrpcOpts.InclusionFunc(parentCtx, method, req, resp) {
			return clientInterceptor(ctx, method, req, resp, cc, invoker, opts...)
		}

		tracer := opentracing.GlobalTracer()

		clientSpan := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
			opentracing.Tag{string(ext.Component), "gRPC"},
		)
		defer clientSpan.Finish()
		ctx = zipkin.InjectSpanContext(ctx, tracer, clientSpan)

		if otgrpcOpts.LogPayloads {
			clientSpan.LogFields(openlog.Object("gRPC request", req))
		}
		err = clientInterceptor(ctx, method, req, resp, cc, invoker, opts...)
		if err == nil {
			if otgrpcOpts.LogPayloads {
				clientSpan.LogFields(openlog.Object("gRPC response", resp))
			}
		} else {
			otgrpc.SetSpanTags(clientSpan, err, true)
			clientSpan.LogFields(openlog.String("event", "error"), openlog.String("message", err.Error()))
		}
		if otgrpcOpts.Decorator != nil {
			otgrpcOpts.Decorator(clientSpan, method, req, resp, err)
		}
		return err
	}
}

func streamClientInterceptorWithZipkin(optFuncs ...zipkin.Option) grpc.StreamClientInterceptor {
	otgrpcOpts := zipkin.NewOptions()
	otgrpcOpts.Apply(optFuncs...)
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		var err error
		var parentCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}
		if otgrpcOpts.InclusionFunc != nil &&
			!otgrpcOpts.InclusionFunc(parentCtx, method, nil, nil) {
			return streamer(ctx, desc, cc, method, opts...)
		}

		tracer := opentracing.GlobalTracer()

		clientSpan := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
			opentracing.Tag{string(ext.Component), "gRPC"},
		)
		ctx = zipkin.InjectSpanContext(ctx, tracer, clientSpan)
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			otgrpc.SetSpanTags(clientSpan, err, true)
			clientSpan.Finish()
			return cs, err
		}
		return newOpenTracingClientStream(cs, method, desc, clientSpan, otgrpcOpts), nil
	}
}

func newOpenTracingClientStream(cs grpc.ClientStream, method string, desc *grpc.StreamDesc, clientSpan opentracing.Span, otgrpcOpts *zipkin.Options) grpc.ClientStream {
	finishChan := make(chan struct{})

	isFinished := new(int32)
	*isFinished = 0
	finishFunc := func(err error) {
		// The current OpenTracing specification forbids finishing a span more than
		// once. Since we have multiple code paths that could concurrently call
		// `finishFunc`, we need to add some sort of synchronization to guard against
		// multiple finishing.
		if !atomic.CompareAndSwapInt32(isFinished, 0, 1) {
			return
		}
		close(finishChan)
		defer clientSpan.Finish()
		if err != nil {
			otgrpc.SetSpanTags(clientSpan, err, true)
		}
		if otgrpcOpts.Decorator != nil {
			otgrpcOpts.Decorator(clientSpan, method, nil, nil, err)
		}
	}
	go func() {
		select {
		case <-finishChan:
			// The client span is being finished by another code path; hence, no
			// action is necessary.
		case <-cs.Context().Done():
			finishFunc(cs.Context().Err())
		}
	}()
	otcs := &zipkin.OpenTracingClientStream{
		ClientStream: cs,
		Desc:         desc,
		FinishFunc:   finishFunc,
	}

	// The `ClientStream` interface allows one to omit calling `Recv` if it's
	// known that the result will be `io.EOF`. See
	// http://stackoverflow.com/q/42915337
	// In such cases, there's nothing that triggers the span to finish. We,
	// therefore, set a finalizer so that the span and the context goroutine will
	// at least be cleaned up when the garbage collector is run.
	runtime.SetFinalizer(otcs, func(otcs *zipkin.OpenTracingClientStream) {
		otcs.FinishFunc(nil)
	})
	return otcs
}

func GetServiceConn(serviceName string) (*grpcpool.ClientConn, error) {
	if poolGet, ok := serviceMap[serviceName]; ok {
		if checkPoolValid(poolGet.pool) {
			conn, err := poolGet.pool.Get(context.Background())
			return conn, err
		} else {
			closePool(serviceName)
		}
	}
	lock.Lock()
	defer lock.Unlock()
	if _, ok := serviceMap[serviceName]; !ok {
		if grpcHost, ok := config.GrpcConfMap.Load(serviceName); ok {

			service := NewServiceGrpcConfig(serviceName, grpcHost.(string))
			poolCreated, poolErr := service.createPoll()
			if poolErr != nil {
				return nil, poolErr
			}
			serviceMap[serviceName] = &pollService{poolCreated}
		} else {
			return nil, errors.New("get grpc host fail")
		}
	}
	conn, err := serviceMap[serviceName].pool.Get(context.Background())
	return conn, err
}

func NewServiceGrpcConfig(name, address string) *ServiceConfig {
	return &ServiceConfig{
		serviceName: name,
		address:     address,
		init:        INIT_NUM,
		capacity:    INIT_CAPACITY,
	}
}

func (service *ServiceConfig) createPoll() (*grpcpool.Pool, error) {
	p, errPoll := grpcpool.New(func() (*grpc.ClientConn, error) {
		return grpc.Dial(service.address, grpc.WithInsecure(), WithClientInterceptor(), WithStreamInterceptor())
	}, service.init, service.capacity, time.Second*5, time.Second*5)
	if errPoll != nil {
		log.Error("init poll fail ", zap.Error(errPoll))
	}
	return p, errPoll
}

func closePool(serviceName string) {
	if _, ok := serviceMap[serviceName]; ok {
		serviceMap[serviceName].pool.Close()
	}
	delete(serviceMap, serviceName)
}

func checkPoolValid(pool *grpcpool.Pool) bool {
	return !pool.IsClosed()
}

func init() {
	if os.Getenv("GRPC_CAPACITY") == "" {
		INIT_CAPACITY = 100
	} else {
		capacityStr := os.Getenv("GRPC_CAPACITY")
		capacity, err := strconv.Atoi(capacityStr)
		if err != nil {
			log.Warn("GRPC_CAPACITY is wrong")
		}
		INIT_CAPACITY = capacity
	}

	config.ListAllDefaultGrpcHost().Range(func(key, value interface{}) bool {
		log.Info("load grpc config", zap.Any("key", key), zap.Any("val", value))
		service := NewServiceGrpcConfig(key.(string), value.(string))
		p, errPool := service.createPoll()
		if errPool != nil {
			log.Warn("get pool error", zap.Error(errPool))
		}

		serviceMap[service.serviceName] = &pollService{p}
		return true
	})

	go func() {
		for {
			rch := site_var.GetDefaultEtcdService().WatchGrpcHostModify()
			for wresp := range rch {
				for _, ev := range wresp.Events {
					if ev.Type == clientv3.EventTypePut {
						splitedList := strings.Split(string(ev.Kv.Key[:]), "/")
						grpcName := splitedList[len(splitedList)-1]

						config.GrpcConfMap.Store(grpcName, string(ev.Kv.Value[:]))

						unRegisterGrpcInstance(grpcName, string(ev.Kv.Value[:]))

						log.Info("update config", zap.ByteString("key", ev.Kv.Key), zap.ByteString("val", ev.Kv.Value))
					}
				}
			}
		}
	}()
}

func unRegisterGrpcInstance(grpcName, renewHost string) {
	lockConfigMap.Lock()
	defer lockConfigMap.Unlock()

	closePool(grpcName)
	grpcConfig := NewServiceGrpcConfig(grpcName, renewHost)
	p, errPool := grpcConfig.createPoll()
	if errPool != nil {
		log.Warn("get pool error", zap.Error(errPool))
	}

	serviceMap[grpcConfig.serviceName] = &pollService{p}
}
