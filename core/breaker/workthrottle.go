package breaker

import "github.com/xsbs1996/go-s-micro/core/logsj"

type workThrottle struct {
	name             string
	internalThrottle internalThrottle
}

// new一个断路器
func newWorkThrottle(name string, t internalThrottle) workThrottle {
	return workThrottle{
		name:             name,
		internalThrottle: t,
	}
}

// allow 断路器手动控制
func (wt workThrottle) allow() (Promise, error) {
	promise, err := wt.internalThrottle.allow()
	return promiseWithReason{
		promise: promise,
	}, err
}

// doReq 断路器自动控制
func (wt workThrottle) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	return wt.internalThrottle.doReq(req, fallback, func(err error) bool {
		return acceptable(err)
	})
}

// promiseWithReason 断路器手动控制结构体
type promiseWithReason struct {
	promise internalPromise
}

// Accept 返回结果正确调用
func (p promiseWithReason) Accept() {
	p.promise.Accept()
}

// Reject 返回结果错误调用
func (p promiseWithReason) Reject(reason string) {
	//输出错误信息
	logsj.BreakerLog(reason)
	p.promise.Reject()
}
