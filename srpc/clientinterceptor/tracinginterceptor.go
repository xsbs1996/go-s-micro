package clientinterceptor

import (
	"context"
	"github.com/xsbs1996/go-s-micro/core/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TracingInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx, span := trace.StartGrpcClientSpan(ctx, cc.Target(), method)
	defer span.Finish()

	var pairs []string
	span.Visit(func(key, val string) bool {
		pairs = append(pairs, key, val)
		return true
	})

	ctx = metadata.AppendToOutgoingContext(ctx, pairs...)

	return invoker(ctx, method, req, reply, cc, opts...)
}
