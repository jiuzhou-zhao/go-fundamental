package servicetoolset

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jiuzhou-zhao/go-fundamental/grpce"
	"github.com/jiuzhou-zhao/go-fundamental/grpce/interceptors"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"google.golang.org/grpc"
)

type ServerToolset struct {
	serverHelper *ServerHelper
	gRpcServer   *grpce.Server
	logger       interfaces.Logger
}

func NewServerToolset(ctx context.Context, logger interfaces.Logger) *ServerToolset {
	if logger == nil {
		logger = &interfaces.EmptyLogger{}
	}
	return &ServerToolset{
		serverHelper: NewServerHelper(ctx, logger),
		logger:       logger,
	}
}

func (st *ServerToolset) Start() {
	if st.gRpcServer != nil {
		st.serverHelper.StartServer(st.gRpcServer)
	}
}

func (st *ServerToolset) Wait() {
	st.serverHelper.Wait()
}

func (st *ServerToolset) CreateGRpcServer(cfg *GRpcConfig, opts []grpc.ServerOption, beforeServerStart func(server *grpc.Server)) error {
	if st.gRpcServer != nil {
		return errors.New("try recreate gRpc server")
	}
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor
	if !cfg.DisableLog {
		if cfg.Logger != nil {
			unaryInterceptors = append(unaryInterceptors, interceptors.ServerLogInterceptor(cfg.Logger))
			streamInterceptors = append(streamInterceptors, interceptors.ServerStreamLogInterceptor(cfg.Logger))
		}
	}

	if cfg.EnableCertVerify {
		if cfg.CertInfo != nil {
			unaryInterceptors = append(unaryInterceptors, interceptors.ServerVerifyInterceptor(cfg.CertInfo))
			streamInterceptors = append(streamInterceptors, interceptors.ServerStreamVerifyInterceptor(cfg.CertInfo))
		}
	}

	if len(unaryInterceptors) > 0 {
		opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)))
	}
	if len(streamInterceptors) > 0 {
		opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)))
	}

	st.gRpcServer = grpce.NewServer(cfg.Address, st.logger, beforeServerStart, opts)
	return nil
}
