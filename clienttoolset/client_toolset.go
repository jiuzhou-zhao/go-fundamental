package clienttoolset

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jiuzhou-zhao/go-fundamental/grpce"
	"github.com/jiuzhou-zhao/go-fundamental/grpce/interceptors"
	"google.golang.org/grpc"
)

func DialGRpcServer(cfg *GRpcClientConfig, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	unaryInterceptors := []grpc.UnaryClientInterceptor{
		interceptors.ClientIDInterceptor(cfg.MetaTransKeys),
	}
	streamInterceptors := []grpc.StreamClientInterceptor{
		interceptors.StreamClientIDInterceptor(cfg.MetaTransKeys),
	}

	if !cfg.DisableLog {
		if cfg.Logger != nil {
			unaryInterceptors = append(unaryInterceptors, interceptors.ClientLogInterceptor(cfg.Logger))
			streamInterceptors = append(streamInterceptors, interceptors.StreamClientLogInterceptor(cfg.Logger))
		}
	}

	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)))
	opts = append(opts, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)))

	return grpce.DialGRpcServer(cfg.Address, &cfg.SecureOption, opts...)
}
