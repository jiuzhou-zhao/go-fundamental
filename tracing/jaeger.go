package tracing

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/grpce/meta"
	"github.com/jiuzhou-zhao/go-fundamental/grpce/utils"
	"github.com/jiuzhou-zhao/go-fundamental/version"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const (
	maxReqResponseMsgLength = 4096
)

var (
	// Morally a const:
	gRPCComponentTag = opentracing.Tag{
		Key:   string(ext.Component),
		Value: "gRPC",
	}
)

func NewTracer(serviceName string, flushInterval time.Duration, address string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: flushInterval,
			LocalAgentHostPort:  address,
		},
	}

	tracer, closer, err = cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
	)
	return
}

func NewGlobalTracer(serviceName string, flushInterval time.Duration, address string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	tracer, closer, err = NewTracer(serviceName, flushInterval, address)
	if err != nil {
		return
	}
	opentracing.SetGlobalTracer(tracer)
	return
}

func ServerOpenTracingInterceptor(tracer opentracing.Tracer, opts ...Option) grpc.UnaryServerInterceptor {
	gRpcOpts := newOptions()
	gRpcOpts.apply(LogPayloads())
	gRpcOpts.apply(opts...)
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		metaData := fmt.Sprintf("metadata %+v", md)
		rw := meta.MDReaderWriter{MD: md}
		spanContext, _ := tracer.Extract(opentracing.TextMap, rw)

		if gRpcOpts.inclusionFunc != nil &&
			!gRpcOpts.inclusionFunc(spanContext, info.FullMethod, req, nil) {
			return handler(ctx, req)
		}

		serverSpan := tracer.StartSpan(
			"server:"+info.FullMethod,
			opentracing.ChildOf(spanContext), // ext.RPCServerOption(spanContext),
			gRPCComponentTag,
			ext.SpanKindRPCServer,
		)

		serverSpan.LogFields(log.String("Metadata", metaData))

		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		if gRpcOpts.logPayloads {
			serverSpan.LogFields(logReqResp("gRPC request", req))
		}
		resp, err = handler(ctx, req)
		if err == nil {
			if gRpcOpts.logPayloads {
				serverSpan.LogFields(logReqResp("gRPC response", resp))
			}
		} else {
			SetSpanTags(serverSpan, err, false)
			serverSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}
		if gRpcOpts.decorator != nil {
			gRpcOpts.decorator(serverSpan, info.FullMethod, req, resp, err)
		}
		serverSpan.Finish()
		return resp, err
	}
}

func ServerStreamOpenTracingInterceptor(tracer opentracing.Tracer, opts ...Option) grpc.StreamServerInterceptor {
	gRpcOpts := newOptions()
	gRpcOpts.apply(opts...)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			md = metadata.New(nil)
		}

		metaData := fmt.Sprintf("metadata %+v", md)
		rw := meta.MDReaderWriter{MD: md}
		spanContext, err := tracer.Extract(opentracing.TextMap, rw)

		if spanContext == nil || err != nil {
			if gRpcOpts.createRootFunc != nil {
				spanContext = gRpcOpts.createRootFunc(info.FullMethod)
			}
			if spanContext == nil {
				parent := tracer.StartSpan(
					info.FullMethod[1:],
					ext.SpanKindRPCClient,
					gRPCComponentTag,
				)
				parent.Finish()
				spanContext = parent.Context()
			}
		}

		if gRpcOpts.inclusionFunc != nil &&
			!gRpcOpts.inclusionFunc(spanContext, info.FullMethod, nil, nil) {
			return handler(srv, ss)
		}

		serverSpan := tracer.StartSpan(
			"server[stream]:"+info.FullMethod,
			opentracing.ChildOf(spanContext), // ext.RPCServerOption(spanContext), opentracing.ChildOf(spanContext)
			opentracing.Tag{
				Key:   "Metadata",
				Value: metaData,
			},
			gRPCComponentTag,
			ext.SpanKindRPCServer,
		)
		serverSpan.LogFields(log.String("metadata", metaData))
		serverSpan.SetTag("app version", version.GetVersionInfo())
		serverSpan.Finish()

		ctx := opentracing.ContextWithSpan(ss.Context(), serverSpan)

		wrapper := utils.NewServerStreamWrapper(ctx, ss)
		err = handler(srv, wrapper)

		serverSpanFinish := tracer.StartSpan(
			"server-finish[stream]:"+info.FullMethod,
			opentracing.ChildOf(serverSpan.Context()),
			gRPCComponentTag,
			ext.SpanKindRPCServer,
		)
		defer serverSpanFinish.Finish()

		if err != nil {
			SetSpanTags(serverSpanFinish, err, false)
			serverSpanFinish.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}
		if gRpcOpts.decorator != nil {
			gRpcOpts.decorator(serverSpan, info.FullMethod, nil, nil, err)
		}
		return err
	}
}

func ClientOpenTracingInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		cliSpan, clientCtx := opentracing.StartSpanFromContextWithTracer(
			ctx,
			tracer,
			"client:"+method,
			gRPCComponentTag,
			ext.SpanKindRPCClient,
		)
		defer cliSpan.Finish()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		mdWriter := meta.MDReaderWriter{MD: md}
		err := tracer.Inject(cliSpan.Context(), opentracing.TextMap, mdWriter)
		if err != nil {
			cliSpan.LogFields(log.String("event", "trace inject error"), log.Error(err))
		}
		ctx = metadata.NewOutgoingContext(clientCtx, md)
		err = invoker(ctx, method, req, resp, cc, opts...)
		if err == nil {
			cliSpan.LogFields(logReqResp("gRPC response", resp))
		} else {
			cliSpan.LogFields(log.String("event", "error"), log.Error(err))
		}

		return err
	}
}

func ClientStreamOpenTracingInterceptor(tracer opentracing.Tracer, opts ...Option) grpc.StreamClientInterceptor {
	gRpcOpts := newOptions()
	gRpcOpts.apply(opts...)
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		var err error
		parent := opentracing.SpanFromContext(ctx)
		if parent == nil {
			parent = tracer.StartSpan(
				method[1:],
				ext.SpanKindRPCClient,
				gRPCComponentTag,
			)
			parent.Finish()
		}
		parentCtx := parent.Context()

		if gRpcOpts.inclusionFunc != nil &&
			!gRpcOpts.inclusionFunc(parentCtx, method, nil, nil) {
			return streamer(ctx, desc, cc, method, opts...)
		}

		clientSpan := tracer.StartSpan(
			"client[stream]:"+method,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
			gRPCComponentTag,
		)
		clientSpan.Finish()

		ctx = InjectSpanContext(ctx, tracer, clientSpan)
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			clientFinishSpan := tracer.StartSpan(
				"client-finish[stream]:"+method,
				opentracing.ChildOf(clientSpan.Context()),
				ext.SpanKindRPCClient,
				gRPCComponentTag,
			)
			clientFinishSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
			SetSpanTags(clientFinishSpan, err, true)
			clientFinishSpan.Finish()
			return cs, err
		}
		return newOpenTracingClientStream(tracer, cs, method, desc, clientSpan, gRpcOpts), nil
	}
}

//
//
//
func logReqResp(key string, reqResp interface{}) log.Field {
	if reqResp == nil {
		return log.Object(key, "nil")
	}
	reqRespProto, ok := reqResp.(proto.Message)
	if ok {
		if ds, err := proto.Marshal(reqRespProto); err == nil {
			log.String(key, truncateSpanMsg(string(ds), maxReqResponseMsgLength))
		}
	}
	return log.Object(key, reqResp)
}

func truncateSpanMsg(msg string, maxLength int) string {
	runeMsg := []rune(msg)
	if len(runeMsg) <= maxLength {
		return msg
	}
	return string(runeMsg[:maxLength-3]) + "..."
}

func InjectSpanContext(ctx context.Context, tracer opentracing.Tracer, clientSpan opentracing.Span) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	mdWriter := meta.MDReaderWriter{MD: md}
	err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, mdWriter)
	// We have no better place to record an error than the Span itself :-/
	if err != nil {
		clientSpan.LogFields(log.String("event", "Tracer.Inject() failed"), log.Error(err))
	}
	return metadata.NewOutgoingContext(ctx, md)
}

func newOpenTracingClientStream(tracer opentracing.Tracer, cs grpc.ClientStream, method string, desc *grpc.StreamDesc, clientSpan opentracing.Span, otgrpcOpts *options) grpc.ClientStream {
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

		clientFinishSpan := tracer.StartSpan(
			"client-finish[stream]:"+method,
			opentracing.ChildOf(clientSpan.Context()),
			ext.SpanKindRPCClient,
			gRPCComponentTag,
		)
		defer clientFinishSpan.Finish()
		if err != nil {
			clientFinishSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
			SetSpanTags(clientFinishSpan, err, true)
		}
		if otgrpcOpts.decorator != nil {
			otgrpcOpts.decorator(clientFinishSpan, method, nil, nil, err)
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

	wrapper := utils.NewClientStreamWrapper(cs, desc, finishFunc)
	// The `ClientStream` interface allows one to omit calling `Recv` if it's
	// known that the result will be `io.EOF`. See
	// http://stackoverflow.com/q/42915337
	// In such cases, there's nothing that triggers the span to finish. We,
	// therefore, set a finalizer so that the span and the context goroutine will
	// at least be cleaned up when the garbage collector is run.
	runtime.SetFinalizer(wrapper, func(wrapper grpc.ClientStream) {
		finishFunc(nil)
	})
	return wrapper
}
