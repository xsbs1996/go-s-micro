package logsj

import (
	"github.com/sirupsen/logrus"
	"github.com/xsbs1996/go-s-micro/core/sys"
)

// BreakerLog 熔断日志
func BreakerLog(reason string) {
	logrus.WithField(timestampKey, GetTimestamp()).
		WithField(callerKey, GetCaller(CallerDepth)).
		WithField(hostName, sys.Hostname()).
		WithField(contentKey, reason).
		Info("Breaker")
}
