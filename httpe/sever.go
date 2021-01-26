package httpe

import (
	"context"
	"net"
	"net/http"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
)

type Server struct {
	address string
	logger  *loge.Logger
	Handler http.Handler
}

func NewServer(address string, logger interfaces.Logger, handler http.Handler) *Server {
	return &Server{
		address: address,
		logger:  loge.NewLogger(logger),
		Handler: handler,
	}
}

func (s *Server) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := &http.Server{
		Handler: s.Handler,
	}

	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return
	}
	s.logger.Infof(ctx, "http server listening on %v", s.address)

	go func() {
		err = server.Serve(l)
		if err != nil {
			s.logger.Errorf(ctx, "http server serve error: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()

	s.logger.Infof(ctx, "http server shutting down")

	_ = server.Close()

	return nil
}
