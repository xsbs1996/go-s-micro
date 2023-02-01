package tracespec

import (
	"context"
	"github.com/gin-gonic/gin"
)

// TracingKey 是上下文的跟踪键
var TracingKey = "X-Trace"

// SetTraceContext 设置链路追踪上下文
func SetTraceContext(gCtx *gin.Context) context.Context {
	ctx := context.Background()
	span, exists := gCtx.Get(TracingKey)
	if !exists {
		return ctx
	}
	return context.WithValue(ctx, TracingKey, span)
}
