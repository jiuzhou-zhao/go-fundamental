package servicetoolset

import (
	"context"
	"errors"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jiuzhou-zhao/go-fundamental/grpce"
	"github.com/jiuzhou-zhao/go-fundamental/grpce/interceptors"
	"github.com/jiuzhou-zhao/go-fundamental/httpe"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"google.golang.org/grpc"
)

type ServerToolset struct {
	ctx          context.Context
	serverHelper *ServerHelper
	gRpcServer   *grpce.Server
	httpServer   *httpe.Server
	logger       interfaces.Logger

	started bool
}

func NewServerToolset(ctx context.Context, logger interfaces.Logger) *ServerToolset {
	if logger == nil {
		logger = &loge.EmptyLogger{}
	}
	return &ServerToolset{
		ctx:          ctx,
		serverHelper: NewServerHelper(ctx, logger),
		logger:       logger,
	}
}

func (st *ServerToolset) Start() error {
	if st.started {
		return errors.New("started")
	}
	st.started = true

	if st.gRpcServer != nil {
		st.serverHelper.StartServer(st.gRpcServer)
	}
	if st.httpServer != nil {
		st.serverHelper.StartServer(st.httpServer)
	}
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
	if cfg == nil || cfg.Address == "" || !strings.Contains(cfg.Address, ":") {
		return errors.New("invalid input args")
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

	st.gRpcServer = grpce.NewServer(cfg.Name, cfg.Address, st.logger, beforeServerStart, opts)
	if cfg.DiscoveryExConfig.Setter != nil {
		st.gRpcServer.EnableDiscovery(cfg.DiscoveryExConfig.Setter, cfg.DiscoveryExConfig.ExternalAddress, cfg.DiscoveryExConfig.Meta)
	}
	return nil
}

func (st *ServerToolset) CreateHttpServer(cfg *HttpServerConfig) error {
	if st.httpServer != nil {
		return errors.New("try recreate http server")
	}
	if cfg == nil || cfg.Address == "" || !strings.Contains(cfg.Address, ":") || cfg.Handler == nil {
		return errors.New("invalid input args")
	}
	st.httpServer = httpe.NewServer(cfg.Address, st.logger, cfg.Handler)

	return nil
}
