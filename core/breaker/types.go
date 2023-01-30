package breaker

// Acceptable 检查调用是否成功
type Acceptable func(err error) bool

// Breaker 断路器接口
type Breaker interface {
	// Name 返回断路器的名称
	Name() string

	// Allow 检查请求是否被允许
	// 如果允许，将返回一个Promise，调用者需要在成功时调用Promise.Accept()，失败时调用 Promise.Reject()
	// 如果不允许, 将返回 ErrServiceUnavailable错误
	Allow() (Promise, error)

	// Do 将运行给定的请求，如果 Breaker 接受请求。
	// Do 立即返回错误，如果 Breaker 拒绝请求。
	// 如果请求中发生恐慌，Breaker 将其作为错误处理并再次导致同样的恐慌。
	Do(req func() error) error

	// DoWithAcceptable 将运行给定的请求，如果 Breaker 接受它。
	// DoWithAcceptable 会立即返回错误，如果 Breaker 拒绝请求。
	// 如果请求中发生恐慌，Breaker 将其作为错误处理并再次导致同样的恐慌。
	// acceptable 检查调用是否成功，即使错误不为零。
	DoWithAcceptable(req func() error, acceptable Acceptable) error

	// DoWithFallback 运行给定的请求，如果 Breaker 接受它。
	// DoWithFallback 运行 fallback 方法，如果 Breaker 拒绝请求。
	// 如果请求中发生恐慌，Breaker 将其作为错误处理并再次导致同样的恐慌。
	DoWithFallback(req func() error, fallback func(err error) error) error

	// DoWithFallbackAcceptable 将运行给定的请求，如果 Breaker 接受它。
	// DoWithFallbackAcceptable 将运行 fallback 方法，如果 Breaker 拒绝请求。
	// 如果请求中发生恐慌，Breaker 将其作为错误处理并再次导致同样的恐慌。
	// acceptable 检查调用是否成功，即使错误不为零。
	DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error
}

// Option 断路器预操作
type Option func(breaker *circuitBreaker)

// internalThrottle 断路器接口
type internalThrottle interface {
	allow() (internalPromise, error)
	doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error
}

// throttle 断路器接口
type throttle interface {
	allow() (Promise, error)
	doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error
}

// Promise 对返回结果进行判断并调用
// 如果返回结果正确，调用 Accept 方法
// 如果返回结果错误，调用 Reject 方法
type Promise interface {
	Accept()
	Reject(reason string)
}

// internalPromise 对返回结果进行判断并调用
// 如果返回结果正确，调用 Accept 方法
// 如果返回结果错误，调用 Reject 方法
type internalPromise interface {
	Accept()
	Reject()
}
