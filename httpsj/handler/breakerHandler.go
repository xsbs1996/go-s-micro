package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xsbs1996/go-s-micro/core/breaker"
	"github.com/xsbs1996/go-s-micro/core/logsj"
	"net/http"
	"strings"
	"sync"
)

type breakerNameList struct {
	rw      sync.RWMutex
	nameMap map[string]breaker.Breaker
}

var breakerNameMap *breakerNameList

func init() {
	breakerNameMap = &breakerNameList{
		rw:      sync.RWMutex{},
		nameMap: make(map[string]breaker.Breaker, 0),
	}
}

func checkName(name string) (breaker.Breaker, bool) {
	breakerNameMap.rw.RLock()
	defer breakerNameMap.rw.RUnlock()
	brk, ok := breakerNameMap.nameMap[name]
	return brk, ok
}

func addName(name string, brk breaker.Breaker) {
	breakerNameMap.rw.Lock()
	defer breakerNameMap.rw.Unlock()
	breakerNameMap.nameMap[name] = brk
}

const breakerSeparator = "://"

// BreakerHandler 熔断中间件
func BreakerHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := strings.Join([]string{ctx.Request.Method, ctx.Request.URL.Path}, breakerSeparator)
		brk, ok := checkName(name)
		if !ok {
			brk = breaker.NewBreaker(breaker.WithName(name))
			addName(name, brk)
		}

		promise, err := brk.Allow()
		if err != nil {
			//第一步判断熔断,无熔断则继续,有熔断直接返回
			logsj.BreakerLog(fmt.Sprintf("%s %s", name, "trigger fuse"))
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
