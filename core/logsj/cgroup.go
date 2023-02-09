package logsj

import (
	"github.com/sirupsen/logrus"
	"github.com/xsbs1996/go-s-micro/core/sys"
)

// CgroupLog Cgroup日志
func CgroupLog(reason error) {
	logrus.WithField(timestampKey, GetTimestamp()).
		WithField(callerKey, GetCaller(CallerDepth)).
		WithField(hostName, sys.Hostname()).
		WithField(contentKey, reason).
		Info("cgroup")
}
