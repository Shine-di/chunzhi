package zipkin

import (
	"context"

	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func InitZipkinConfig(serviceName string) {
	// env := config.GetEnv()
	// path := "/" + env + "/config/zipkin/host"

	// url, errUrl := site_var.GetDefaultEtcdService().Get(path)
	// if errUrl != nil {
	// 	log.Warn("", zap.Error(errUrl))
	// }

	// // create collector.
	// collector, err := zipkintracer.NewHTTPCollector(url)
	// if err != nil {
	// 	log.Warn("", zap.Error(err))
	// }

	// // create recorder.
	// recorder := zipkintracer.NewRecorder(collector, true, getIp(), serviceName)

	// // create tracer.
	// tracer, err := zipkintracer.NewTracer(
	// 	recorder,
	// 	zipkintracer.ClientServerSameSpan(true),
	// 	zipkintracer.TraceID128Bit(true),
	// )
	// if err != nil {
	// 	log.Warn("", zap.Error(err))
	// }
	// // explicitly set our tracer to be the default tracer.
	// opentracing.SetGlobalTracer(tracer)
}

func ZipkinWrapper(ot opentracing.Tracer) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx, span, err := traceIntoContext(ctx, ot, req.Endpoint())
			if err != nil {
				return err
			}
			defer span.Finish()
			return h(ctx, req, rsp)
		}
	}
}

func traceIntoContext(ctx context.Context, tracer opentracing.Tracer, name string) (context.Context, opentracing.Span, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	// copy the metadata to prevent race
	md = metadata.Copy(md)

	var sp opentracing.Span
	wireContext, err := tracer.Extract(opentracing.HTTPHeaders, MetadataMircoReaderWriter{md})
	if err != nil {
		sp = tracer.StartSpan(name)
	} else {
		sp = tracer.StartSpan(name, ext.RPCServerOption(wireContext))
	}
	if err := sp.Tracer().Inject(sp.Context(), opentracing.HTTPHeaders, MetadataMircoReaderWriter{md}); err != nil {
		return nil, nil, err
	}
	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)

	return ctx, sp, nil
}
