package logsj

import (
	"github.com/sirupsen/logrus"
	"os"
)

const (
	callerKey    = "caller"
	contentKey   = "content"
	hostName     = "hostname"
	timestampKey = "timestamp"
)

func init() {
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
	})
}
