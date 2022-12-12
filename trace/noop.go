package trace

import (
	"github.com/gin-gonic/gin"
	"github.com/xsbs1996/go-s-micro/trace/tracespec"
)

var emptyNoopSpan = noopSpan{}

type noopSpan struct{}

func (s noopSpan) Finish() {
}

func (s noopSpan) Follow(ctx *gin.Context, serviceName, operationName string) tracespec.Trace {
	return emptyNoopSpan
}

func (s noopSpan) Fork(ctx *gin.Context, serviceName, operationName string) tracespec.Trace {
	return emptyNoopSpan
}

func (s noopSpan) SpanID() string {
	return ""
}

func (s noopSpan) TraceID() string {
	return ""
}

func (s noopSpan) Visit(fn func(key, val string) bool) {
}
