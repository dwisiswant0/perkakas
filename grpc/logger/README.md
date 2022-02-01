# Logger grpc interceptor
This package contains Unary and Streaming grpc interceptor for Logger.
Please refer to instructions below for datails usage

## Unary Server Interceptor

```go
import(
    ...
    "github.com/kitabisa/perkakas/grpcinterceptor/logger"
    ...
)

func main(){
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 50051))
	if err != nil {
		opt.Logger.Error(fmt.Sprintf("failed to listen %s:%d", host, port), err, nil)
	}

    // use the interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			logger.UnaryServerInterceptor,
		),
	)

    // initialize grpc handler
	grpcHandler := grpcHandler.NewFlagHandler(opt)

	pb.RegisterFlagServer(grpcServer, grpcHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	opt.Logger.Info(fmt.Sprintf("GRPC serve at %s:%d", host, port), nil)

	grpcServer.Serve(lis)
}
```

### Using Unary Server Interceptor with custom RequestID context key

```go
import(
    ...
    "github.com/kitabisa/perkakas/grpcinterceptor/logger"
    ...
)

func main(){
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 50051))
	if err != nil {
		opt.Logger.Error(fmt.Sprintf("failed to listen %s:%d", host, port), err, nil)
	}

    // initialize logger interceptor 
    // with requestid context key options
    interceptor := logger.NewInterceptor(
		WithRequestIDContextKey("custom-requestid-context-key"),
	)
    
    // use the interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.UnaryServerInterceptor,
		),
	)

    // initialize grpc handler
	grpcHandler := grpcHandler.NewFlagHandler(opt)

	pb.RegisterFlagServer(grpcServer, grpcHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	opt.Logger.Info(fmt.Sprintf("GRPC serve at %s:%d", host, port), nil)

	grpcServer.Serve(lis)
}
```

## Streaming Server Interceptor

```go
import(
    ...
    "github.com/kitabisa/perkakas/grpcinterceptor/logger"
    ...
)

func main(){
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 50051))
	if err != nil {
		opt.Logger.Error(fmt.Sprintf("failed to listen %s:%d", host, port), err, nil)
	}

    // use the interceptor
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(
			logger.StreamingServerInterceptor,
		),
	)

    // initialize grpc handler
	grpcHandler := grpcHandler.NewFlagHandler(opt)

	pb.RegisterFlagServer(grpcServer, grpcHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	opt.Logger.Info(fmt.Sprintf("GRPC serve at %s:%d", host, port), nil)

	grpcServer.Serve(lis)
}
```

### Using Streaming Server Interceptor with custom RequestID context key

```go
import(
    ...
    "github.com/kitabisa/perkakas/grpcinterceptor/logger"
    ...
)

func main(){
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 50051))
	if err != nil {
		opt.Logger.Error(fmt.Sprintf("failed to listen %s:%d", host, port), err, nil)
	}

    // initialize logger interceptor 
    // with requestid context key options
    interceptor := logger.NewInterceptor(
		WithRequestIDContextKey("custom-requestid-context-key"),
	)
    
    // use the interceptor
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(
			interceptor.StreamingServerInterceptor,
		),
	)

    // initialize grpc handler
	grpcHandler := grpcHandler.NewFlagHandler(opt)

	pb.RegisterFlagServer(grpcServer, grpcHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	opt.Logger.Info(fmt.Sprintf("GRPC serve at %s:%d", host, port), nil)

	grpcServer.Serve(lis)
}
```

## Using both Unary and Streaming server interceptor

```go
import(
    ...
    "github.com/kitabisa/perkakas/grpcinterceptor/logger"
    ...
)

func main(){
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 50051))
	if err != nil {
		opt.Logger.Error(fmt.Sprintf("failed to listen %s:%d", host, port), err, nil)
	}

    // use the interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			logger.UnaryServerInterceptor,
		),
        grpc.StreamInterceptor(
			logger.StreamingServerInterceptor,
		),
	)

    // initialize grpc handler
	grpcHandler := grpcHandler.NewFlagHandler(opt)

	pb.RegisterFlagServer(grpcServer, grpcHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	opt.Logger.Info(fmt.Sprintf("GRPC serve at %s:%d", host, port), nil)

	grpcServer.Serve(lis)
}
```