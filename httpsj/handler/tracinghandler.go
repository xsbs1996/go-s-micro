package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xsbs1996/go-s-micro/core/sys"
	"github.com/xsbs1996/go-s-micro/core/trace"
)

// TracingHandler 链路追踪中间件
func TracingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		carrier, err := trace.Extract(trace.HttpFormat, ctx.Request.Header)
		if err != nil && err != trace.ErrInvalidCarrier {
			logrus.Error(err)
		}

		span := trace.StartServerSpan(ctx, carrier, sys.Hostname(), ctx.Request.RequestURI)
		defer span.Finish()

		ctx.Next()
	}

}
