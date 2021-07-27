package requestid

import (
	"context"
	"testing"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/kitabisa/perkakas/v2/grpcinterceptor/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryInterceptorWithoutRequestID(t *testing.T) {

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.NotEmpty(t, requestID)

		return requestID, nil
	}

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{})
	UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestInstanceUnaryInterceptorWithoutRequestID(t *testing.T) {

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.NotEmpty(t, requestID)

		return requestID, nil
	}

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

	interceptor := NewInterceptor()
	interceptor.UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)

}

func TestUnaryInterceptorWithRequestID(t *testing.T) {
	reqID := uuid.NewV4().String()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID)

		assert.Equal(t, reqID, requestID)

		return requestID, nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestInstanceUnaryInterceptorWithRequestID(t *testing.T) {
	reqID := uuid.NewV4().String()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID)

		assert.Equal(t, reqID, requestID)

		return requestID, nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	interceptor := NewInterceptor()
	interceptor.UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestInstanceUnaryInterceptorWithRequestIDAndCustomMetadataKey(t *testing.T) {
	reqID := uuid.NewV4().String()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID)

		assert.Equal(t, reqID, requestID)

		return requestID, nil
	}

	ctx := context.Background()
	md := metadata.Pairs("custom-key", reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	interceptor := NewInterceptor(
		WithMetadataKey("custom-key"),
	)
	interceptor.UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestInstanceUnaryInterceptorWithRequestIDAndCustomContextKey(t *testing.T) {
	reqID := uuid.NewV4().String()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		requestID := ctx.Value(ctxkeys.ContextKey("custom-key"))

		assert.Equal(t, reqID, requestID)

		return requestID, nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	interceptor := NewInterceptor(
		WithContextKey("custom-key"),
	)
	interceptor.UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestStreamingServerInterceptorWithoutRequestID(t *testing.T) {
	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.NotEmpty(t, requestID)

		return nil
	}

	serverStream := mocks.NewMockServerStream()
	StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestStreamingServerInterceptorWithRequestID(t *testing.T) {
	reqID := uuid.NewV4().String()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.Equal(t, reqID, requestID)

		return nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	serverStream := mocks.NewMockServerStream()
	serverStream.SetContext(ctx)

	StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithoutRequestID(t *testing.T) {
	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.NotEmpty(t, requestID)

		return nil
	}

	serverStream := mocks.NewMockServerStream()

	interceptor := NewInterceptor()
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithRequestID(t *testing.T) {
	reqID := uuid.NewV4().String()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.Equal(t, reqID, requestID)

		return nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	serverStream := mocks.NewMockServerStream()
	serverStream.SetContext(ctx)

	interceptor := NewInterceptor()
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithRequestIDAndCustomMetadataKey(t *testing.T) {
	reqID := uuid.NewV4().String()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.Equal(t, reqID, requestID)

		return nil
	}

	ctx := context.Background()
	md := metadata.Pairs("custom-key", reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	serverStream := mocks.NewMockServerStream()
	serverStream.SetContext(ctx)

	interceptor := NewInterceptor(
		WithMetadataKey("custom-key"),
	)
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithRequestIDAndCustomContextKey(t *testing.T) {
	reqID := uuid.NewV4().String()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.ContextKey("custom-key")).(string)
		assert.Equal(t, reqID, requestID)

		return nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	serverStream := mocks.NewMockServerStream()
	serverStream.SetContext(ctx)

	interceptor := NewInterceptor(
		WithContextKey("custom-key"),
	)
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}
