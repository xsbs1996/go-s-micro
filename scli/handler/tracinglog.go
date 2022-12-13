package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xsbs1996/go-s-micro/logsj"
	"github.com/xsbs1996/go-s-micro/trace"
	"github.com/xsbs1996/go-s-micro/trace/tracespec"
	"io"
	"net/http"
	"sync"
)

const defaultMemory = 32 << 20

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

		logsj.TracingLog(span.Operation(), span.SpanID(), span.TraceID(), requestInputs(ctx))
		ctx.Next()
	}
}

// RequestInputs 获取所有参数
func requestInputs(ctx *gin.Context) map[string]interface{} {
	contentType := ctx.ContentType()

	var (
		dataMap  = make(map[string]interface{})
		queryMap = make(map[string]interface{})
		postMap  = make(map[string]interface{})
	)

	for k := range ctx.Request.URL.Query() {
		queryMap[k] = ctx.Query(k)
	}

	switch contentType {
	case "application/json":
		var bodyBytes []byte
		if ctx.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(ctx.Request.Body)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if ctx.Request != nil && ctx.Request.Body != nil {
			if err := json.NewDecoder(ctx.Request.Body).Decode(&postMap); err != nil {
				return nil
			}
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	case "multipart/form-data":
		if err := ctx.Request.ParseMultipartForm(defaultMemory); err != nil {
			return nil
		}
		for k, v := range ctx.Request.PostForm {
			if len(v) > 1 {
				postMap[k] = v
			} else if len(v) == 1 {
				postMap[k] = v[0]
			}
		}
	default:
		if err := ctx.Request.ParseForm(); err != nil {
			return nil
		}
		if err := ctx.Request.ParseMultipartForm(defaultMemory); err != nil {
			if err != http.ErrNotMultipart {
				return nil
			}
		}
		for k, v := range ctx.Request.PostForm {
			if len(v) > 1 {
				postMap[k] = v
			} else if len(v) == 1 {
				postMap[k] = v[0]
			}
		}

	}

	var mu sync.RWMutex
	for k, v := range queryMap {
		mu.Lock()
		dataMap[k] = v
		mu.Unlock()
	}
	for k, v := range postMap {
		mu.Lock()
		dataMap[k] = v
		mu.Unlock()
	}

	return dataMap
}
