package interceptors

import (
	"context"
	"google.golang.org/grpc"
)

type streamWrapper struct {
	grpc.ServerStream
	WrapperContext context.Context
}

func (s *streamWrapper) Context() context.Context {
	return s.WrapperContext
}

func (s *streamWrapper) RecvMsg(m interface{}) error {
	return s.ServerStream.RecvMsg(m)
}

func (s *streamWrapper) SendMsg(m interface{}) error {
	return s.ServerStream.SendMsg(m)
}
