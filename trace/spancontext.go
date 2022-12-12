package trace

type spanContext struct {
	traceID string // 链路ID(全追踪链路唯一)
	spanID  string // spanID(在本节点上ID,向下文传递时同时传递此ID作为下文parentSpanId)
}

// TraceID 获取traceID
func (sc spanContext) TraceID() string {
	return sc.traceID
}

// SpanID 获取spanID
func (sc spanContext) SpanID() string {
	return sc.spanID
}

func (sc spanContext) Visit(fn func(key, val string) bool) {
	fn(traceIdKey, sc.traceID)
	fn(spanIdKey, sc.spanID)
}
