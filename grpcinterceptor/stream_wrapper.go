package grpcinterceptor

import (
	"context"

	"google.golang.org/grpc"
)

type ServerStreamWrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func NewServerStreamWrapper(ctx context.Context, stream grpc.ServerStream) *ServerStreamWrapper {
	return &ServerStreamWrapper{
		ctx:          ctx,
		ServerStream: stream,
	}
}

func (s *ServerStreamWrapper) Context() context.Context {
	return s.ctx
}
