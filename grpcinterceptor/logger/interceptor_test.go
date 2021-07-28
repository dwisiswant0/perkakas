package logger

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/kitabisa/perkakas/v2/grpcinterceptor/mocks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnaryInterceptorWithoutRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.NotContains(t, out.String(), ctxkeys.CtxXKtbsRequestID)

		return nil, nil
	}

	ctx := context.Background()
	UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestUnaryInterceptorWithRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	reqID := uuid.NewV4().String()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.Contains(t, out.String(), reqID)
		assert.Contains(t, out.String(), ctxkeys.CtxXKtbsRequestID)

		return nil, nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkeys.CtxXKtbsRequestID, reqID)
	UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestInstanceUnaryInterceptorWithoutRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.NotContains(t, out.String(), ctxkeys.CtxXKtbsRequestID)

		return nil, nil
	}

	ctx := context.Background()

	interceptor := NewInterceptor()
	interceptor.UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestInstanceUnaryInterceptorWithRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	reqID := uuid.NewV4().String()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.Contains(t, out.String(), reqID)
		assert.Contains(t, out.String(), ctxkeys.CtxXKtbsRequestID)

		return nil, nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkeys.CtxXKtbsRequestID, reqID)

	interceptor := NewInterceptor()
	interceptor.UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestInstanceUnaryInterceptorWithRequestIDAndCustomContextKey(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	reqID := uuid.NewV4().String()
	reqIDContextKey := ctxkeys.ContextKey("custom-key")

	test := func(ctx context.Context, req interface{}) (interface{}, error) {

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.Contains(t, out.String(), reqID)
		assert.Contains(t, out.String(), reqIDContextKey)

		return nil, nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, reqIDContextKey, reqID)

	interceptor := NewInterceptor(
		WithRequestIDContextKey("custom-key"),
	)

	interceptor.UnaryServerInterceptor(ctx, nil, mocks.UnaryInfo, test)
}

func TestStreamingServerInterceptorWithoutRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.NotContains(t, out.String(), ctxkeys.CtxXKtbsRequestID)

		return nil
	}

	serverStream := mocks.NewMockServerStream(false)
	StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestStreamingServerInterceptorWithRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	reqID := uuid.NewV4().String()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.Contains(t, out.String(), reqID)
		assert.Contains(t, out.String(), ctxkeys.CtxXKtbsRequestID)

		return nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkeys.CtxXKtbsRequestID, reqID)

	serverStream := mocks.NewMockServerStream(false)
	serverStream.SetContext(ctx)

	StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithoutRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.NotContains(t, out.String(), ctxkeys.CtxXKtbsRequestID)

		return nil
	}

	serverStream := mocks.NewMockServerStream(false)

	interceptor := NewInterceptor()
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	reqID := uuid.NewV4().String()

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.Contains(t, out.String(), reqID)
		assert.Contains(t, out.String(), ctxkeys.CtxXKtbsRequestID)

		return nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkeys.CtxXKtbsRequestID, reqID)

	serverStream := mocks.NewMockServerStream(false)
	serverStream.SetContext(ctx)

	interceptor := NewInterceptor()
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}

func TestInstanceStreamingServerInterceptorWithRequestIDAndCustomContextKey(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	reqID := uuid.NewV4().String()

	reqIDContextKey := ctxkeys.ContextKey("custom-key")

	test := func(srv interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.Contains(t, out.String(), reqID)
		assert.Contains(t, out.String(), reqIDContextKey)

		return nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, reqIDContextKey, reqID)

	serverStream := mocks.NewMockServerStream(false)
	serverStream.SetContext(ctx)

	interceptor := NewInterceptor(
		WithRequestIDContextKey("custom-key"),
	)
	interceptor.StreamingServerInterceptor(nil, serverStream, mocks.StreamInfo, test)
}
