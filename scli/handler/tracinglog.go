package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xsbs1996/go-s-micro/logsj"
	"github.com/xsbs1996/go-s-micro/trace"
	"github.com/xsbs1996/go-s-micro/trace/tracespec"
)

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
		logsj.TracingLog(span, ctx.Request)
		ctx.Next()
	}
}
