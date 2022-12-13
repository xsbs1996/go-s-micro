package tracespec

type SpanContext interface {
	TraceID() string
	SpanID() string
	ParentSpanID() string
	Visit(fn func(key, val string) bool)
}
