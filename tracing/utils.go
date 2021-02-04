package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func TryTracing(ctx context.Context, handler func(span opentracing.Span)) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		handler(span)
	}
}
