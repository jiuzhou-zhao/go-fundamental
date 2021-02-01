package tracing

import "github.com/opentracing/opentracing-go"

type Option func(o *options)

func LogPayloads() Option {
	return func(o *options) {
		o.logPayloads = true
	}
}

type SpanInclusionFunc func(
	parentSpanCtx opentracing.SpanContext,
	method string,
	req, resp interface{}) bool

func IncludingSpans(inclusionFunc SpanInclusionFunc) Option {
	return func(o *options) {
		o.inclusionFunc = inclusionFunc
	}
}

type SpanDecoratorFunc func(
	span opentracing.Span,
	method string,
	req, resp interface{},
	grpcError error)

func SpanDecorator(decorator SpanDecoratorFunc) Option {
	return func(o *options) {
		o.decorator = decorator
	}
}

type SpanCreateRootFunc func(
	method string) opentracing.SpanContext

func CreateRootSpans(createRootFunc SpanCreateRootFunc) Option {
	return func(o *options) {
		o.createRootFunc = createRootFunc
	}
}

type options struct {
	logPayloads bool
	// May be nil.
	decorator SpanDecoratorFunc
	// May be nil.
	inclusionFunc SpanInclusionFunc
	// May be nil.
	createRootFunc SpanCreateRootFunc
}

// newOptions returns the default options.
func newOptions() *options {
	return &options{
		logPayloads: false,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
