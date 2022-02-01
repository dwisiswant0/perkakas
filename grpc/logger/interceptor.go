package logger

import (
	"context"
	"github.com/kitabisa/perkakas/v3/perkakas-grpc/wrapper"
	"sync"

	"github.com/kitabisa/perkakas/v3/perkakas-grpc/ctxkeys"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	instance *Interceptor
	doOnce   sync.Once
)

// Init creating default interceptor instance
// the reason of not using go init() is to prevent
// unwanted extra logger interceptor instance when using this interceptor
// without default instance
func Init() {
	doOnce.Do(func() {
		instance = NewInterceptor()
	})
}

type Options func(*Interceptor)

type Interceptor struct {
	requestIDContextKey ctxkeys.ContextKey
}

// WithRequestIDContextKey set requestID context value key.
// provides an option to use this interceptor with
// context value key other than "X-Ktbs-Request-ID"
func WithRequestIDContextKey(key string) Options {
	return func(i *Interceptor) {
		i.requestIDContextKey = ctxkeys.ContextKey(key)
	}
}

// UnaryServerInterceptor calling logger UnaryServerInterceptor
// with default interceptor instance
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	Init()
	return instance.UnaryServerInterceptor(ctx, req, info, handler)
}

// StreamingServerInterceptor calling logger StreamingServerInterceptor
// with default interceptor instance
func StreamingServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	Init()
	return instance.StreamingServerInterceptor(srv, stream, info, handler)
}

func NewInterceptor(opts ...Options) *Interceptor {
	i := &Interceptor{}

	for _, opt := range opts {
		opt(i)
	}

	if i.requestIDContextKey == "" {
		i.requestIDContextKey = ctxkeys.CtxXKtbsRequestID
	}

	return i
}

func (i *Interceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	logger := log.Logger

	reqID, ok := ctx.Value(i.requestIDContextKey).(string)
	if ok {
		logger = log.With().Str(i.requestIDContextKey.String(), reqID).Logger()
	}

	ctx = context.WithValue(ctx, ctxkeys.CtxLogger, logger)

	resp, err = handler(ctx, req)

	return
}

func (i *Interceptor) StreamingServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := stream.Context()
	logger := log.Logger

	reqID, ok := ctx.Value(i.requestIDContextKey).(string)
	if ok {
		logger = log.With().Str(i.requestIDContextKey.String(), reqID).Logger()
	}

	ctx = context.WithValue(ctx, ctxkeys.CtxLogger, logger)
	newStream := wrapper.NewServerStreamWrapper(ctx, stream)

	return handler(srv, newStream)
}
