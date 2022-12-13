package serverinterceptors

import (
	"context"
	"github.com/xsbs1996/go-s-micro/logsj"
	"github.com/xsbs1996/go-s-micro/trace"
	"github.com/xsbs1996/go-s-micro/trace/tracespec"
	"google.golang.org/grpc"
)

// TracingLog grpc输出链路追踪拦截器
func TracingLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if span, ok := ctx.Value(tracespec.TracingKey).(*trace.Span); ok {
			logsj.TracingLog(span.Operation(), span.SpanID(), span.TraceID(), req)
			return handler(ctx, req)
		}
		return handler(ctx, req)
	}
}
