package tracespec

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Trace interface {
	SpanContext
	Finish()
	Fork(ctx *gin.Context, serviceName, operationName string) Trace
	GrpcFork(ctx context.Context, serviceName, operationName string) (context.Context, Trace)
	Follow(ctx *gin.Context, serviceName, operationName string) Trace
}
