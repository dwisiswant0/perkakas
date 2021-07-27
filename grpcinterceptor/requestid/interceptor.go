package requestid

import (
	"context"
	"errors"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/kitabisa/perkakas/v2/grpcinterceptor"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type reqIDcontextKey string

const (
	GrpcRequestIDKey = "x-ktbs-request-id"
)

type Options func(*Interceptor)

type Interceptor struct {
	reqIDKey reqIDcontextKey
}

func WithReqIDKey(reqIDKey string) Options {
	return func(i *Interceptor) {
		i.reqIDKey = reqIDcontextKey(reqIDKey)
	}
}

func NewInterceptor(opts ...Options) *Interceptor {
	i := &Interceptor{}

	for _, opt := range opts {
		opt(i)
	}

	if i.reqIDKey == "" {
		i.reqIDKey = reqIDcontextKey(ctxkeys.CtxXKtbsRequestID)
	}

	return i
}

// UnaryServerInterceptor get interceptor without creating Interceptor instance
// TODO: TO BE DEPRECATED
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	interceptor := NewInterceptor()
	return interceptor.UnaryServerInterceptor(ctx, req, info, handler)
}

func (i *Interceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	reqID, err := getRequestID(ctx, GrpcRequestIDKey)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if reqID == "" {
		reqID = uuid.NewV4().String()
	}

	ctx = context.WithValue(ctx, ctxkeys.CtxXKtbsRequestID, reqID)

	resp, err = handler(ctx, req)

	return
}

func (i *Interceptor) StreamingServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := stream.Context()
	reqID, err := getRequestID(ctx, GrpcRequestIDKey)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	if reqID == "" {
		reqID = uuid.NewV4().String()
	}

	ctx = context.WithValue(ctx, i.reqIDKey, reqID)
	newStream := grpcinterceptor.NewServerStreamWrapper(ctx, stream)

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
