package trace

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xsbs1996/go-s-micro/trace/tracespec"
	"github.com/xsbs1996/go-s-micro/utils/stringfunc"
	"github.com/xsbs1996/go-s-micro/utils/timefunc"
	"strconv"
	"strings"
	"time"
)

const (
	initSpanID  = "0"      // 初始span
	clientFlag  = "client" // 客户端span
	serverFlag  = "server" // 服务端span
	spanSepRune = '.'      // 分隔符
)

var spanSep = string([]byte{spanSepRune})

type Span struct {
	ctx           spanContext // span上下文
	serviceName   string      // 服务名
	operationName string      // 操作
	startTime     time.Time   // 开始时间
	flag          string      // 操作标记
	children      int         // 本span fork出来的children
}

// newServerSpan 初始span
func newServerSpan(carrier Carrier, serviceName string, operationName string) tracespec.Trace {
	traceId := stringfunc.TakeWithPriority(func() string {
		if carrier != nil {
			return carrier.Get(traceIdKey)
		}
		return ""
	}, stringfunc.RandId)

	spanId := stringfunc.TakeWithPriority(func() string {
		if carrier != nil {
			return carrier.Get(spanIdKey)
		}
		return ""
	}, func() string {
		return initSpanID
	})

	return &Span{
		ctx: spanContext{
			traceID: traceId,
			spanID:  spanId,
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     timefunc.Time(),
		flag:          serverFlag,
	}

}

func (s *Span) Finish() {}

func (s *Span) Follow(ctx *gin.Context, serviceName, operationName string) tracespec.Trace {
	span := &Span{
		ctx: spanContext{
			traceID: s.ctx.traceID,
			spanID:  s.followSpanID(),
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     timefunc.Time(),
		flag:          s.flag,
	}
	ctx.Set(tracespec.TracingKey, span)
	return span
}

// Fork 分支上下文跟踪
func (s *Span) Fork(ctx *gin.Context, serviceName, operationName string) tracespec.Trace {
	span := &Span{
		ctx: spanContext{
			traceID: s.ctx.traceID,
			spanID:  s.forkSpanID(),
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     timefunc.Time(),
		flag:          clientFlag,
	}
	ctx.Set(tracespec.TracingKey, span)
	return span
}

func (s *Span) SpanID() string {
	return s.ctx.SpanID()
}

func (s *Span) TraceID() string {
	return s.ctx.TraceID()
}

func (s *Span) Visit(fn func(key, val string) bool) {
	s.ctx.Visit(fn)
}

// 生成fork的子spanID
func (s *Span) forkSpanID() string {
	s.children++
	return fmt.Sprintf("%s.%d", s.ctx.spanID, s.children)
}

// 生成follow的子spanID
func (s *Span) followSpanID() string {
	fields := strings.FieldsFunc(s.ctx.spanID, func(r rune) bool {
		return r == spanSepRune
	})
	if len(fields) == 0 {
		return s.ctx.spanID
	}

	last := fields[len(fields)-1]
	val, err := strconv.Atoi(last)
	if err != nil {
		return s.ctx.spanID
	}

	last = strconv.Itoa(val + 1)
	fields[len(fields)-1] = last

	return strings.Join(fields, spanSep)
}

func StartClientSpan(ctx *gin.Context, serviceName, operationName string) tracespec.Trace {
	spanI, exists := ctx.Get(tracespec.TracingKey)
	if !exists {
		return emptyNoopSpan
	}

	if span, ok := spanI.(*Span); ok {
		return span.Fork(ctx, serviceName, operationName)
	}
	return emptyNoopSpan
}

func StartServerSpan(ctx *gin.Context, carrier Carrier, serviceName, operationName string) tracespec.Trace {
	span := newServerSpan(carrier, serviceName, operationName)
	ctx.Set(tracespec.TracingKey, span)
	return span
}