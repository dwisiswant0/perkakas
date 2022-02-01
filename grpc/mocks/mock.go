package mocks

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func NewMockServerStream(withMetada bool) *MockServerStream {
	ctx := context.Background()

	if withMetada {
		ctx = metadata.NewIncomingContext(ctx, metadata.MD{})
	}

	return &MockServerStream{
		ctx: ctx,
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
