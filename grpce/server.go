package grpce

import (
	"context"
	"net"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"google.golang.org/grpc"
)

type BeforeServerStart func(server *grpc.Server)

type Server struct {
	address           string
	logger            *loge.Logger
	beforeServerStart BeforeServerStart
	opts              []grpc.ServerOption
}

func NewServer(address string, logger interfaces.Logger, beforeServerStart BeforeServerStart, opts []grpc.ServerOption) *Server {
	return &Server{
		address:           address,
		logger:            loge.NewLogger(logger),
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
	s.logger.Infof(ctx, "grpc server listening on %v", s.address)

	go func() {
		err = server.Serve(l)
		if err != nil {
			s.logger.Errorf(ctx, "grpc server serve error: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()

	s.logger.Infof(ctx, "grpc server shutting down")

	server.Stop()

	return nil
}
