package breaker

import "sync"

var (
	lock     sync.RWMutex
	breakers = make(map[string]Breaker)
)

// GetBreaker 通过读写锁获取breaker
func GetBreaker(name string) Breaker {
	lock.RLock()
	b, ok := breakers[name]
	lock.RUnlock()
	if ok {
		return b
	}

	lock.Lock()
	b, ok = breakers[name]
	if !ok {
		b = NewBreaker(WithName(name))
		breakers[name] = b
	}
	lock.Unlock()

	return b
}

// 以下方法参考Breaker的注释,作用是一样的
func do(name string, execute func(b Breaker) error) error {
	return execute(GetBreaker(name))
}

func Do(name string, req func() error) error {
	return do(name, func(b Breaker) error {
		return b.Do(req)
	})
}

func DoWithAcceptable(name string, req func() error, acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithAcceptable(req, acceptable)
	})
}

func DoWithFallback(name string, req func() error, fallback func(err error) error) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallback(req, fallback)
	})
}

func DoWithFallbackAcceptable(name string, req func() error, fallback func(err error) error,
	acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallbackAcceptable(req, fallback, acceptable)
	})
}
