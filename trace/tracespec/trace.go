package tracespec

import (
	"github.com/gin-gonic/gin"
)

type Trace interface {
	SpanContext
	Finish()
	Fork(ctx *gin.Context, serviceName, operationName string) Trace
	Follow(ctx *gin.Context, serviceName, operationName string) Trace
}
