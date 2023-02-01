package serverinterceptors

import (
	"context"
	"github.com/xsbs1996/go-s-micro/core/logsj"
	"github.com/xsbs1996/go-s-micro/core/trace"
	"github.com/xsbs1996/go-s-micro/core/trace/tracespec"
	"google.golang.org/grpc"
)

// TracingLog grpc输出链路追踪拦截器
func TracingLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if span, ok := ctx.Value(tracespec.TracingKey).(*trace.Span); ok {
			defer logsj.TracingLog(span.Operation(), span.SpanID(), span.TraceID(), span.StartTime(), req)
			return handler(ctx, req)
		}
		return handler(ctx, req)
	}
}
