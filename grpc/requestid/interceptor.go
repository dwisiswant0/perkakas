package requestid

import (
	"context"
	"errors"
	"github.com/kitabisa/perkakas/perkakas-grpc/wrapper"
	"sync"

	"github.com/kitabisa/perkakas/perkakas-grpc/ctxkeys"
	uuid "github.com/kitabisa/perkakas/random"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	GrpcRequestIDKey = "x-ktbs-request-id"
)

var (
	instance *Interceptor
	doOnce   sync.Once
)

// Init creating default interceptor instance
// the reason of not using go init() is to prevent
// unwanted extra requestID interceptor instance when using this interceptor
// without default instance
func Init() {
	doOnce.Do(func() {
		instance = NewInterceptor()
	})
}

// UnaryServerInterceptor calling requestID UnaryServerInterceptor
// with default interceptor instance
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	Init()
	return instance.UnaryServerInterceptor(ctx, req, info, handler)
}

// StreamingServerInterceptor calling requestID StreamingServerInterceptor
// with default interceptor instance
func StreamingServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	Init()
	return instance.StreamingServerInterceptor(srv, stream, info, handler)
}

type Options func(*Interceptor)

type Interceptor struct {
	metadataKey string
	contextKey  ctxkeys.ContextKey
}

// WithMetadataKey set reqIDKey that should be send by client
// using grpc metadata, provides an option to use this interceptor with
// requestIDKey metadata other than "x-ktbs-request-id"
func WithMetadataKey(key string) Options {
	return func(i *Interceptor) {
		i.metadataKey = key
	}
}

// WithContextKey set requestID context value key.
// provides an option to use this interceptor with
// context value key other than "X-Ktbs-Request-ID"
func WithContextKey(key string) Options {
	return func(i *Interceptor) {
		i.contextKey = ctxkeys.ContextKey(key)
	}
}

func NewInterceptor(opts ...Options) *Interceptor {
	i := &Interceptor{}

	for _, opt := range opts {
		opt(i)
	}

	if i.metadataKey == "" {
		i.metadataKey = GrpcRequestIDKey
	}

	if i.contextKey == "" {
		i.contextKey = ctxkeys.CtxXKtbsRequestID
	}

	return i
}

func (i *Interceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	reqID, err := getRequestID(ctx, i.metadataKey)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if reqID == "" {
		reqID = uuid.UUID()
	}

	ctx = context.WithValue(ctx, i.contextKey, reqID)

	resp, err = handler(ctx, req)

	return
}

func (i *Interceptor) StreamingServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := stream.Context()
	reqID, err := getRequestID(ctx, i.metadataKey)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	if reqID == "" {
		reqID = uuid.UUID()
	}

	ctx = context.WithValue(ctx, i.contextKey, reqID)
	newStream := wrapper.NewServerStreamWrapper(ctx, stream)

	return handler(srv, newStream)
}

func getRequestID(ctx context.Context, key string) (val string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("failed to retrieve metadata")
	}

	v := md[key]

	if len(v) == 0 {
		return
	}

	val = v[0]

	return
}
