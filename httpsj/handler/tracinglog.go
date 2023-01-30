package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xsbs1996/go-s-micro/core/logsj"
	"github.com/xsbs1996/go-s-micro/trace"
	"github.com/xsbs1996/go-s-micro/trace/tracespec"
	"github.com/xsbs1996/go-s-micro/utils/ginfunc"
)

// TracingLog 链路追踪日志中间件
func TracingLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanI, exists := ctx.Get(tracespec.TracingKey)
		if !exists {
			ctx.Next()
			return
		}
		span, ok := spanI.(*trace.Span)
		if !ok {
			ctx.Next()
			return
		}

		defer logsj.TracingLog(span.Operation(), span.SpanID(), span.TraceID(), span.StartTime(), ginfunc.RequestInputs(ctx))
		ctx.Next()
	}
}
