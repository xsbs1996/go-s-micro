package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xsbs1996/go-s-micro/core/breaker"
	"net/http"
	"strings"
)

const breakerSeparator = " "

// BreakerHandler 熔断中间件
func BreakerHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := strings.Join([]string{ctx.Request.Method, ctx.Request.URL.Path}, breakerSeparator)
		brk := breaker.GetBreaker(name)

		promise, err := brk.Allow()
		if err != nil {
			//第一步判断熔断,无熔断则继续,有熔断直接返回
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, "")
			return
		}

		//defer写入到滑动时间窗口
		defer func() {
			//如果http code 小于500,则结果正确,否则结果错误
			if ctx.Writer.Status() < http.StatusInternalServerError {
				//增加正确请求数量
				promise.Accept()
			} else {
				//增加总请求数量
				promise.Reject(fmt.Sprintf("%d %s", ctx.Writer.Status(), http.StatusText(ctx.Writer.Status())))
			}

		}()
		ctx.Next()
	}

}
