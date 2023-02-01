package trace

type spanContext struct {
	traceID      string // 链路ID(全追踪链路唯一)
	spanID       string // spanID
	parentSpanID string
}

// TraceID 获取traceID
func (sc spanContext) TraceID() string {
	return sc.traceID
}

// SpanID 获取spanID
func (sc spanContext) SpanID() string {
	return sc.spanID
}

// ParentSpanID 获取spanID
func (sc spanContext) ParentSpanID() string {
	return sc.parentSpanID
}

func (sc spanContext) Visit(fn func(key, val string) bool) {
	fn(traceIdKey, sc.traceID)
	fn(spanIdKey, sc.spanID)
}
