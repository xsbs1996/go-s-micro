package tracespec

type SpanContext interface {
	TraceID() string
	SpanID() string
	Visit(fn func(key, val string) bool)
}
