//go:build !linux

package cpustat

// CpuCurrentRate 返回cpu目前占用率
func CpuCurrentRate() float64 {
	return 0
}
