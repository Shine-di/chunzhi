package zipkin

import "github.com/opentracing/opentracing-go"

type Option func(o *Options)

func LogPayloads() Option {
	return func(o *Options) {
		o.LogPayloads = true
	}
}

type SpanInclusionFunc func(
	parentSpanCtx opentracing.SpanContext,
	method string,
	req, resp interface{}) bool

func IncludingSpans(inclusionFunc SpanInclusionFunc) Option {
	return func(o *Options) {
		o.InclusionFunc = inclusionFunc
	}
}

type SpanDecoratorFunc func(
	span opentracing.Span,
	method string,
	req, resp interface{},
	grpcError error)

func SpanDecorator(decorator SpanDecoratorFunc) Option {
	return func(o *Options) {
		o.Decorator = decorator
	}
}

type Options struct {
	LogPayloads   bool
	Decorator     SpanDecoratorFunc
	InclusionFunc SpanInclusionFunc
}

func NewOptions() *Options {
	return &Options{
		LogPayloads:   false,
		InclusionFunc: nil,
	}
}

func (o *Options) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}
