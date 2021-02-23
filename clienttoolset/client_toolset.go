package clienttoolset

import (
	"context"
	"errors"
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jiuzhou-zhao/go-fundamental/grpce"
	"github.com/jiuzhou-zhao/go-fundamental/grpce/interceptors"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/jiuzhou-zhao/go-fundamental/tracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

const (
	rrGRpcServerConfig = `
{
	"loadBalancingConfig": [ { "round_robin": {} } ]
}
`
)

func DialGRpcServer(cfg *GRpcClientConfig, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	if cfg == nil {
		return nil, errors.New("no cfg")
	}
	unaryInterceptors := []grpc.UnaryClientInterceptor{
		interceptors.ClientIDInterceptor(cfg.MetaTransKeys),
	}
	streamInterceptors := []grpc.StreamClientInterceptor{
		interceptors.ClientStreamIDInterceptor(cfg.MetaTransKeys),
	}

	if !cfg.DisableLog {
		if cfg.Logger != nil {
			unaryInterceptors = append(unaryInterceptors, interceptors.ClientLogInterceptor(cfg.Logger))
			streamInterceptors = append(streamInterceptors, interceptors.ClientStreamLogInterceptor(cfg.Logger))
		}
	}

	if !cfg.DisableLog {
		if cfg.Logger != nil {
			unaryInterceptors = append(unaryInterceptors, interceptors.ClientLogInterceptor(cfg.Logger))
			streamInterceptors = append(streamInterceptors, interceptors.ClientStreamLogInterceptor(cfg.Logger))
		}
	}
	if cfg.EnableTracing {
		tracingObj := opentracing.GlobalTracer()

		if cfg.TracingConfig.ServerAddr != "" {
			var err error
			tracingObj, _, err = tracing.NewTracer(cfg.TracingConfig.ServiceName, cfg.TracingConfig.FlushInterval, cfg.TracingConfig.ServerAddr)
			if err != nil {
				return nil, fmt.Errorf("new tracker failed: %v", err)
			}
		}

		if tracingObj != nil {
			unaryInterceptors = append(unaryInterceptors, tracing.ClientOpenTracingInterceptor(tracingObj))
			streamInterceptors = append(streamInterceptors, tracing.ClientStreamOpenTracingInterceptor(tracingObj))
		}
	}

	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)))
	opts = append(opts, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)))

	return grpce.DialGRpcServer(cfg.Address, &cfg.SecureOption, opts...)
}

func DialGRpcServerByName(schema, serverName string, cfg *GRpcClientConfig, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithDefaultServiceConfig(rrGRpcServerConfig))
	if cfg == nil {
		cfg = &GRpcClientConfig{}
	}
	cfg.Address = fmt.Sprintf("%s:///%s", schema, serverName)
	return DialGRpcServer(cfg, opts)
}

func RegisterSchemas(ctx context.Context, cfg *RegisterSchemasConfig) error {
	if cfg == nil {
		return errors.New("no cfg")
	}
	if cfg.Getter == nil {
		return errors.New("no getter")
	}

	logger := loge.NewLogger(cfg.Logger)

	for _, schema := range cfg.Schemas {
		err := grpce.RegisterResolver(cfg.Getter, cfg.Logger, schema)
		if err != nil {
			logger.Errorf(ctx, "register schema %v failed: %v", schema, err)
		}
	}
	return nil
}
