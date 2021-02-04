package tracing

import (
	"context"
	"fmt"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type Logger struct {
}

func NewTracingLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Record(ctx context.Context, depth int, level interfaces.LoggerLevel, v ...interface{}) {
	TryTracing(ctx, func(span opentracing.Span) {
		span.LogFields(log.String(level.String(), truncateSpanMsg(fmt.Sprint(v...), maxReqResponseMsgLength)))
	})
}

func (logger *Logger) Recordf(ctx context.Context, depth int, level interfaces.LoggerLevel, format string, v ...interface{}) {
	TryTracing(ctx, func(span opentracing.Span) {
		span.LogFields(log.String(level.String(), truncateSpanMsg(fmt.Sprintf(format, v...), maxReqResponseMsgLength)))
	})
}
