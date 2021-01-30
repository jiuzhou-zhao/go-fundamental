package clienttoolset

import (
	"context"
	"errors"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jiuzhou-zhao/go-fundamental/grpce"
	"github.com/jiuzhou-zhao/go-fundamental/grpce/interceptors"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
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

func DialGRpcServerByName(schema, serverName string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	return DialGRpcServer(&GRpcClientConfig{
		Address: fmt.Sprintf("%s:///%s", schema, serverName),
	}, []grpc.DialOption{grpc.WithDefaultServiceConfig(rrGRpcServerConfig)})
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
