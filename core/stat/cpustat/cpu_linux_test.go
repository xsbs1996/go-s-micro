package cpustat

import (
	"fmt"
	"testing"
	"time"
)

func TestCpuCurrentRate(t *testing.T) {
	tt := time.NewTicker(time.Second * 1)
	defer tt.Stop()

	for {
		select {
		case <-tt.C:
			fmt.Println("cup占用率", RefreshCpu())
		}
	}
}
