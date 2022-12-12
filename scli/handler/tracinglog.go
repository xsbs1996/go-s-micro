package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xsbs1996/go-s-micro/trace"
	"github.com/xsbs1996/go-s-micro/trace/tracespec"
	"github.com/xsbs1996/go-s-micro/utils/logfunc"
)

const (
	callerKey    = "caller"
	contentKey   = "content"
	operationKey = "operation"
	spanKey      = "span"
	timestampKey = "@timestamp"
	traceKey     = "trace"
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

		logrus.WithField(timestampKey, logfunc.GetTimestamp()).
			WithField(callerKey, logfunc.GetCaller(logfunc.CallerDepth)).
			WithField(operationKey, span.ServiceOperation()).
			WithField(spanKey, span.SpanID()).
			WithField(traceKey, span.TraceID()).
			WithField(contentKey, span).Info("Tracing")

		ctx.Next()
	}

}
