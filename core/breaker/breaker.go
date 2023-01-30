package breaker

import (
	"errors"
	"github.com/xsbs1996/go-s-micro/utils/stringfunc"
)

// ErrServiceUnavailable 当断路器状态打开时返回
var ErrServiceUnavailable = errors.New("circuit breaker is open")

// 断路器结构体
type circuitBreaker struct {
	name     string // 短路器名称
	throttle        // 代理，短路器业务全部由 throttle 实现
}

// NewBreaker new一个断路器
func NewBreaker(opts ...Option) Breaker {
	var b circuitBreaker
	for _, opt := range opts {
		opt(&b)
	}

	if len(b.name) == 0 {
		b.name = stringfunc.Rand()
	}

	b.throttle = newWorkThrottle(b.name, newGoogleBreaker())

	return &b
}

func (cb *circuitBreaker) Name() string {
	return cb.name
}

func (cb *circuitBreaker) Allow() (Promise, error) {
	return cb.throttle.allow()
}

func (cb *circuitBreaker) Do(req func() error) error {
	return cb.throttle.doReq(req, nil, defaultAcceptable)
}

func (cb *circuitBreaker) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	return cb.throttle.doReq(req, nil, acceptable)
}

func (cb *circuitBreaker) DoWithFallback(req func() error, fallback func(err error) error) error {
	return cb.throttle.doReq(req, fallback, defaultAcceptable)
}

func (cb *circuitBreaker) DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	return cb.throttle.doReq(req, fallback, acceptable)
}

// 断路器默认验证方法
func defaultAcceptable(err error) bool {
	return err == nil
}

// WithName 断路器设置名称
func WithName(name string) Option {
	return func(b *circuitBreaker) {
		b.name = name
	}
}
