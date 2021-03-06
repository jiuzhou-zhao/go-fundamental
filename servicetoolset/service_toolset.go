package servicetoolset

import (
	"context"
	"errors"
	"fmt"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jiuzhou-zhao/go-fundamental/grpce"
	"github.com/jiuzhou-zhao/go-fundamental/grpce/interceptors"
	"github.com/jiuzhou-zhao/go-fundamental/httpe"
	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/jiuzhou-zhao/go-fundamental/tracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type ServerToolset struct {
	ctx          context.Context
	logger       *loge.Logger
	serverHelper *ServerHelper
	gRpcServer   *grpce.Server
	httpServer   *httpe.Server

	started bool
}

func NewServerToolset(ctx context.Context, logger interfaces.Logger) *ServerToolset {
	sst := &ServerToolset{
		ctx:    ctx,
		logger: loge.NewLogger(logger),
	}
	sst.serverHelper = NewServerHelper(ctx, sst.logger.GetLogger())
	return sst
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
		interceptors.ServerStreamIDInterceptor(cfg.MetaTransKeys),
	}
	if !cfg.DisableLog {
		logger := cfg.Logger
		if logger == nil {
			logger = st.logger.GetLogger()
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

	if cfg.EnableTracing {
		tracingObj := opentracing.GlobalTracer()

		if cfg.TracingConfig.ServerAddr != "" {
			var err error
			tracingObj, _, err = tracing.NewTracer(cfg.TracingConfig.ServiceName, cfg.TracingConfig.FlushInterval, cfg.TracingConfig.ServerAddr)
			if err != nil {
				return fmt.Errorf("new tracker failed: %v", err)
			}
		}

		if tracingObj == nil {
			return errors.New("no valid tracing config")
		}
		unaryInterceptors = append(unaryInterceptors, tracing.ServerOpenTracingInterceptor(tracingObj))
		streamInterceptors = append(streamInterceptors, tracing.ServerStreamOpenTracingInterceptor(tracingObj))
	}

	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)))
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)))

	st.gRpcServer = grpce.NewServer(cfg.Name, cfg.Address, st.logger.GetLogger(), beforeServerStart, opts)
	if cfg.DiscoveryExConfig.Setter != nil {
		st.gRpcServer.EnableDiscovery(cfg.DiscoveryExConfig.Setter, cfg.DiscoveryExConfig.ExternalAddress, cfg.DiscoveryExConfig.Meta)
	}
	if cfg.EnableGRpcWeb {
		if cfg.GRpcWebAddress == "" {
			st.logger.Fatal(st.ctx, "gRpcWeb server enabled, but no address config")
		}
		st.gRpcServer.EnableGRpcWeb(cfg.GRpcWebAddress, cfg.GRpcWebUseWebsocket, cfg.GRpcWebPingInterval)
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
	st.httpServer = httpe.NewServer(cfg.Address, st.logger.GetLogger(), cfg.Handler)

	return nil
}
