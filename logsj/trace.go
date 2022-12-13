package logsj

import (
	"github.com/sirupsen/logrus"
	"github.com/xsbs1996/go-s-micro/utils/logfunc"
	"github.com/xsbs1996/go-s-micro/utils/sysfunc"
)

const (
	callerKey    = "caller"
	contentKey   = "content"
	hostName     = "hostname"
	operationKey = "operation"
	spanKey      = "span"
	timestampKey = "timestamp"
	traceKey     = "trace"
)

func TracingLog(operation, spanID, traceID, req interface{}) {
	logrus.WithField(timestampKey, logfunc.GetTimestamp()).
		WithField(callerKey, logfunc.GetCaller(logfunc.CallerDepth)).
		WithField(hostName, sysfunc.Hostname()).
		WithField(operationKey, operation).
		WithField(spanKey, spanID).
		WithField(traceKey, traceID).
		WithField(contentKey, req).Info("Tracing")
}
