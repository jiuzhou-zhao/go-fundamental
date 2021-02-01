package interceptors

import (
	"context"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/fmtutils"
	"github.com/jiuzhou-zhao/go-fundamental/grpce/meta"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"google.golang.org/grpc"
)

func ServerLogInterceptor(log interfaces.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		id := meta.IdFromOutgoingContext(ctx)

		eLog := loge.NewLogger(log)
		eLog.Infof(ctx, "[SRV][REQ] id:%v method:%v req:\n%v",
			id, info.FullMethod, fmtutils.Marshal(req))

		st := time.Now()

		res, err := handler(ctx, req)

		eLog.Infof(ctx, "[SRV][RESP] id:%v method:%v cost: %v err:%v data:\n%v;]",
			id, info.FullMethod, time.Since(st), err, fmtutils.Marshal(res))

		return res, err
	}
}

func ServerStreamLogInterceptor(log interfaces.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		id := meta.IdFromOutgoingContext(ss.Context())

		eLog := loge.NewLogger(log)
		eLog.Infof(ss.Context(), "[SRV][REQ][STREAM] id:%v method:%v connected", id, info.FullMethod)

		st := time.Now()

		err := handler(srv, ss)

		eLog.Infof(ss.Context(), "[SRV][RESP][STREAM] id:%v method:%v closed. cost:%v", id, info.FullMethod, time.Since(st))

		return err
	}
}

func ClientLogInterceptor(log interfaces.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		id := meta.IdFromOutgoingContext(ctx)

		eLog := loge.NewLogger(log)
		eLog.Infof(ctx, "[CLI][REQ] id:%v method:%v target:%v req:\n%v",
			id, method, cc.Target(), fmtutils.Marshal(req))

		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		eLog.Infof(ctx, "[CLI][RESP] id:%v method:%v cost:%v err:%v data:\n%v",
			id, method, time.Since(start), err, fmtutils.Marshal(reply))

		return err
	}
}

func ClientStreamLogInterceptor(log interfaces.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer,
		opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
		id := meta.IdFromOutgoingContext(ctx)

		eLog := loge.NewLogger(log)
		eLog.Infof(ctx, "[CLI][STREAM] id:%v method:%v connected",
			id, method)

		return streamer(ctx, desc, cc, method, opts...)
	}
}
