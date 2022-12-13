package serverinterceptors

import (
	"context"
	"fmt"
	"github.com/xsbs1996/go-s-micro/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryTracingInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		carrier, err := trace.Extract(trace.GrpcFormat, md)
		if err != nil {
			return handler(ctx, req)
		}

		ctx, span := trace.StartGrpcServerSpan(ctx, carrier, serviceName, info.FullMethod)
		fmt.Println(span)
		defer span.Finish()
		return handler(ctx, req)
	}
}
