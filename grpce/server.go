package grpce

import (
	"context"
	"net"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"google.golang.org/grpc"
)

type BeforeServerStart func(server *grpc.Server)

type Server struct {
	address           string
	logger            interfaces.Logger
	beforeServerStart BeforeServerStart
	opts              []grpc.ServerOption
}

func NewServer(address string, logger interfaces.Logger, beforeServerStart BeforeServerStart, opts []grpc.ServerOption) *Server {
	if logger == nil {
		logger = &interfaces.EmptyLogger{}
	}
	return &Server{
		address:           address,
		logger:            logger,
		beforeServerStart: beforeServerStart,
		opts:              opts,
	}
}

func (s *Server) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := grpc.NewServer(s.opts...)

	if s.beforeServerStart != nil {
		s.beforeServerStart(server)
	}

	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return
	}
	s.logger.Recordf(ctx, interfaces.LogLevelInfo, "gRpcServer listening on %v", s.address)

	go func() {
		err = server.Serve(l)
		if err != nil {
			s.logger.Recordf(ctx, interfaces.LogLevelError, "gRpcServer serve error: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()

	s.logger.Recordf(ctx, interfaces.LogLevelInfo, "gRpcServer shutting down")

	server.Stop()

	return nil
}
