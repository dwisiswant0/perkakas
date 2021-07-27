package mocks

import (
	"context"

	"google.golang.org/grpc"
)

var UnaryInfo = &grpc.UnaryServerInfo{
	FullMethod: "TestUnaryInterceptor",
}

var StreamInfo = &grpc.StreamServerInfo{
	FullMethod:     "TestServerInterceptor",
	IsServerStream: true,
}

type MockServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func NewMockServerStream() *MockServerStream {
	return &MockServerStream{
		ctx: context.Background(),
	}
}

func (m *MockServerStream) SetContext(ctx context.Context) {
	m.ctx = ctx
}

func (m *MockServerStream) Context() context.Context {
	return m.ctx
}

func (m *MockServerStream) SendMsg(msg interface{}) error {
	return nil
}

func (m *MockServerStream) RecvMsg(msg interface{}) error {
	return nil
}
