package clientinterceptor

import (
	"context"
	"github.com/xsbs1996/go-s-micro/core/breaker"
	"github.com/xsbs1996/go-s-micro/srpc/codes"
	"google.golang.org/grpc"
	"path"
)

func BreakerInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	breakerName := path.Join(cc.Target(), method)

	//rpc处理方法
	processing := func() error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	return breaker.DoWithAcceptable(breakerName, processing, codes.Acceptable)
}
