package interceptors

import (
	"context"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/fmtutils"
	"github.com/jiuzhou-zhao/go-fundamental/grpce/meta"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"google.golang.org/grpc"
)

func ServerLogInterceptor(log interfaces.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		id := meta.IdFromIncomingContext(ctx)

		log.Infof("[req] id:%v method:%v req:\n%v",
			id, info.FullMethod, fmtutils.Marshal(req))

		st := time.Now()

		ctx = meta.IdToOutgoingContext(ctx, id)

		res, err := handler(ctx, req)

		log.Infof("[rsp] id:%v method:%v cost:%v err:%v data:\n%v",
			id, info.FullMethod, time.Since(st), err, fmtutils.Marshal(res))

		return res, err
	}
}

func ServerStreamLogInterceptor(log interfaces.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		id := meta.IdFromIncomingContext(ss.Context())

		wrapper := &streamWrapper{
			ServerStream:   ss,
			WrapperContext: meta.IdToOutgoingContext(ss.Context(), id),
		}

		log.Infof("[stream] id:%v method:%v connected", id, info.FullMethod)

		st := time.Now()

		err := handler(srv, wrapper)

		log.Infof("[stream] id:%v method:%v closed. time:%v", id, info.FullMethod, time.Since(st))

		return err
	}
}
