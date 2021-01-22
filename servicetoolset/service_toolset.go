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
	ctx          context.Context
	serverHelper *ServerHelper
	gRpcServer   *grpce.Server
	logger       interfaces.Logger
}

func NewServerToolset(ctx context.Context, logger interfaces.Logger) *ServerToolset {
	if logger == nil {
		logger = &interfaces.EmptyLogger{}
	}
	return &ServerToolset{
		ctx:          ctx,
		serverHelper: NewServerHelper(ctx, logger),
		logger:       logger,
	}
}

func (st *ServerToolset) Start() error {
	if st.gRpcServer != nil {
		return errors.New("started")
	}
	st.serverHelper.StartServer(st.gRpcServer)
	return nil
}

func (st *ServerToolset) Wait() {
	_ = st.Start()
	st.serverHelper.Wait()
}

func (st *ServerToolset) CreateGRpcServer(cfg *GRpcServerConfig, opts []grpc.ServerOption, beforeServerStart func(server *grpc.Server)) error {
	if st.gRpcServer != nil {
		return errors.New("try recreate gRpc server")
	}
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		interceptors.ServerIDInterceptor(cfg.MetaTransKeys),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		interceptors.StreamServerIDInterceptor(cfg.MetaTransKeys),
	}
	if !cfg.DisableLog {
		logger := cfg.Logger
		if logger == nil {
			logger = st.logger
		}
		if logger != nil {
			unaryInterceptors = append(unaryInterceptors, interceptors.ServerLogInterceptor(logger))
			streamInterceptors = append(streamInterceptors, interceptors.ServerStreamLogInterceptor(logger))
		}
	}

	if cfg.EnableCertVerify {
		if cfg.CertInfo != nil {
			unaryInterceptors = append(unaryInterceptors, interceptors.ServerVerifyInterceptor(cfg.CertInfo))
			streamInterceptors = append(streamInterceptors, interceptors.ServerStreamVerifyInterceptor(cfg.CertInfo))
		}
	}

	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)))
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)))

	st.gRpcServer = grpce.NewServer(cfg.Address, st.logger, beforeServerStart, opts)

	return nil
}
