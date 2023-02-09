package logsj

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xsbs1996/go-s-micro/core/sys"
	"time"
)

const (
	operationKey = "operation"
	spanKey      = "span"
	traceKey     = "trace"
	runTime      = "runTime"
)

// TracingLog 链路追踪日志
func TracingLog(operation, spanID, traceID interface{}, startTime time.Time, req interface{}) {
	logrus.WithField(timestampKey, GetTimestamp()).
		WithField(callerKey, GetCaller(CallerDepth)).
		WithField(hostName, sys.Hostname()).
		WithField(operationKey, operation).
		WithField(spanKey, spanID).
		WithField(traceKey, traceID).
		WithField(runTime, fmt.Sprintf("%d%s", time.Since(startTime).Milliseconds(), "ms")).
		WithField(contentKey, req).
		Info("Tracing")
}
