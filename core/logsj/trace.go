package logsj

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xsbs1996/go-s-micro/utils/logfunc"
	"github.com/xsbs1996/go-s-micro/utils/sysfunc"
	"time"
)

const (
	callerKey    = "caller"
	contentKey   = "content"
	hostName     = "hostname"
	operationKey = "operation"
	spanKey      = "span"
	timestampKey = "timestamp"
	traceKey     = "trace"
	runtime      = "runtime"
)

// TracingLog 链路追踪日志
func TracingLog(operation, spanID, traceID interface{}, startTime time.Time, req interface{}) {
	logrus.WithField(timestampKey, logfunc.GetTimestamp()).
		WithField(callerKey, logfunc.GetCaller(logfunc.CallerDepth)).
		WithField(hostName, sysfunc.Hostname()).
		WithField(operationKey, operation).
		WithField(spanKey, spanID).
		WithField(traceKey, traceID).
		WithField(runtime, fmt.Sprintf("%d%s", time.Since(startTime).Milliseconds(), "ms")).
		WithField(contentKey, req).
		Info("Tracing")
}

// BreakerLog 熔断日志
func BreakerLog(reason string) {
	logrus.WithField(hostName, sysfunc.Hostname()).WithField(contentKey, reason).Info("Breaker")
}
