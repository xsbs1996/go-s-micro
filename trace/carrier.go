package trace

import (
	"errors"
	"net/http"
	"strings"
)

var ErrInvalidCarrier = errors.New("invalid carrier")

type Carrier interface {
	Get(key string) string
	Set(key, value string)
}

type httpCarrier http.Header         // http载荷
type grpcCarrier map[string][]string // grpc载荷(key不区分大小写)

// Get 从http请求头中获取数据
func (h httpCarrier) Get(key string) string {
	return http.Header(h).Get(key)
}

// Set 向http请求头中插入数据
func (h httpCarrier) Set(key string, value string) {
	http.Header(h).Set(key, value)
}

func (g grpcCarrier) Get(key string) string {
	if val, ok := g[strings.ToLower(key)]; ok && len(val) > 0 {
		return val[0]
	} else {
		return ""
	}
}

func (g grpcCarrier) Set(key, val string) {
	key = strings.ToLower(key)
	g[key] = append(g[key], val)
}
