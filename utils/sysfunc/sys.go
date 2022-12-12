package sysfunc

import (
	"github.com/xsbs1996/go-s-micro/utils/stringfunc"
	"os"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = stringfunc.RandId()
	}
}

func Hostname() string {
	return hostname
}
