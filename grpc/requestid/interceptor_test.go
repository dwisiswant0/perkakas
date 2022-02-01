package requestid

import (
	"context"
	"github.com/kitabisa/perkakas/perkakas-grpc/mocks"
	"testing"

	"github.com/kitabisa/perkakas/perkakas-grpc/ctxkeys"
	uuid "github.com/kitabisa/perkakas/random"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryInterceptorWithoutGrpcMetadata(t *testing.T) {

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	ctx := context.Background()
	_, err := UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
	assert.NotEmpty(t, err)
}

func TestInstanceUnaryInterceptorWithoutGrpcMetadata(t *testing.T) {

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	ctx := context.Background()

	interceptor := NewInterceptor()
	_, err := interceptor.UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
	assert.NotEmpty(t, err)
}

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
	reqID := uuid.UUID()

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
	reqID := uuid.UUID()

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
	reqID := uuid.UUID()

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
	reqID := uuid.UUID()

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

func TestStreamingInterceptorWithoutGrpcMetadata(t *testing.T) {

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		return nil
	}

	ctx := context.Background()

	serverStream := mocks.NewMockServerStream(true)
	serverStream.SetContext(ctx)

	err := StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)

	assert.NotEmpty(t, err)
}

func TestStreamingServerInterceptorWithoutRequestID(t *testing.T) {
	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.NotEmpty(t, requestID)

		return nil
	}

	serverStream := mocks.NewMockServerStream(true)
	StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestStreamingServerInterceptorWithRequestID(t *testing.T) {
	reqID := uuid.UUID()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.Equal(t, reqID, requestID)

		return nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	serverStream := mocks.NewMockServerStream(true)
	serverStream.SetContext(ctx)

	StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingInterceptorWithoutGrpcMetadata(t *testing.T) {

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		return nil
	}

	ctx := context.Background()

	serverStream := mocks.NewMockServerStream(true)
	serverStream.SetContext(ctx)

	interceptor := NewInterceptor()
	err := interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)

	assert.NotEmpty(t, err)
}

func TestInstanceStreamingServerInterceptorWithoutRequestID(t *testing.T) {
	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.NotEmpty(t, requestID)

		return nil
	}

	serverStream := mocks.NewMockServerStream(true)

	interceptor := NewInterceptor()
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithRequestID(t *testing.T) {
	reqID := uuid.UUID()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.Equal(t, reqID, requestID)

		return nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	serverStream := mocks.NewMockServerStream(true)
	serverStream.SetContext(ctx)

	interceptor := NewInterceptor()
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithRequestIDAndCustomMetadataKey(t *testing.T) {
	reqID := uuid.UUID()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.Equal(t, reqID, requestID)

		return nil
	}

	ctx := context.Background()
	md := metadata.Pairs("custom-key", reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	serverStream := mocks.NewMockServerStream(true)
	serverStream.SetContext(ctx)

	interceptor := NewInterceptor(
		WithMetadataKey("custom-key"),
	)
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithRequestIDAndCustomContextKey(t *testing.T) {
	reqID := uuid.UUID()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		requestID := ctx.Value(ctxkeys.ContextKey("custom-key")).(string)
		assert.Equal(t, reqID, requestID)

		return nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)

	serverStream := mocks.NewMockServerStream(true)
	serverStream.SetContext(ctx)

	interceptor := NewInterceptor(
		WithContextKey("custom-key"),
	)
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}
