package cpustat

import (
	"errors"
	"fmt"
	"github.com/xsbs1996/go-s-micro/core/logsj"
	"github.com/xsbs1996/go-s-micro/core/stat/cpustat/cgroup"
	"github.com/xsbs1996/go-s-micro/utils/iofunc"
	"github.com/xsbs1996/go-s-micro/utils/strconvfunc"
	"strings"
	"time"
)

const (
	cpuTicks  = 100 // /proc/stat下的CPU时间单位为 jiffies， 1 jiffies = 0.01秒
	cpuFields = 8   // cpu统计数据格式
)

var (
	preSystem uint64 // cpu实际使用时间
	quota     uint64 // 实际使用的cpu核心数量
	preTotal  uint64 // cgroup中所有task的cpu使用时长
	cores     uint64 // cgroup中cpu核心总数量
)

func init() {
	cpus, err := cgroup.GetCpuacctUsagePercpu()
	if err != nil {
		logsj.CgroupLog(err)
		return
	}
	cores = uint64(len(cpus)) // cpu核心数总数

	sets, err := cgroup.GetCpusetCpus()
	if err != nil {
		logsj.CgroupLog(err)
		return
	}
	quota = uint64(len(sets))            // 实际使用的cpu核心数量
	cq, err := cgroup.GetCpuCfsQuotaUs() // 通过cgroup设定的周期时间，周期最大可用时间进行验证可用cpu核心数量
	if err == nil {
		if cq > 0 {
			period, err := cgroup.GetCpuCfsPeriodUs()
			if err != nil {
				logsj.CgroupLog(err)
				return
			}
			if period > 0 {
				// 计算实际可用的核心数量
				limit := uint64(cq / period)
				if limit < quota {
					quota = limit
				}
			}

		}
	}

	preSystem, err = systemCpuUsage() // 计算CPU自系统开启以来的各项指标的累计时间
	if err != nil {
		logsj.CgroupLog(err)
		return
	}

	preTotal, err = cgroup.GetCpuacctUsage() // 计算cgroup中所有task的cpu使用时长
	if err != nil {
		logsj.CgroupLog(err)
		return
	}
}

// RefreshCpu 返回cpu目前占用率
func RefreshCpu() uint64 {
	total, err := cgroup.GetCpuacctUsage()
	if err != nil {
		return 0
	}
	system, err := systemCpuUsage()
	if err != nil {
		return 0
	}

	var usage uint64

	// 计算差值
	cpuDelta := total - preTotal
	systemDelta := system - preSystem

	// 计算占用率
	if cpuDelta > 0 && systemDelta > 0 {
		usage = uint64((float64(cpuDelta) * float64(cores)) / (float64(systemDelta) * float64(quota)) * 1e3)
	}

	// 更新值
	preSystem = system
	preTotal = total

	return usage
}

// 计算CPU自系统开启以来的各项指标的累计时间,实际cpu时间
// user	    从系统启动开始累计到当前时刻，用户态的CPU时间（单位：jiffies） ，不包含 nice值为负进程。1jiffies=0.01秒
// nice	    从系统启动开始累计到当前时刻，nice值为负的进程所占用的CPU时间（单位：jiffies）
// system	从系统启动开始累计到当前时刻，核心时间（单位：jiffies）
// idle	    从系统启动开始累计到当前时刻，除硬盘IO等待时间以外其它空闲时间（单位：jiffies）
// iowait	从系统启动开始累计到当前时刻，硬盘IO等待时间（单位：jiffies）
// irq	    从系统启动开始累计到当前时刻，硬中断时间（单位：jiffies）
// softirq	从系统启动开始累计到当前时刻，软中断时间（单位：jiffies）
func systemCpuUsage() (uint64, error) {
	lines, err := iofunc.ReadTextLines("/proc/stat", iofunc.WithoutBlank())
	if err != nil {
		return 0, err
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			if len(fields) < cpuFields {
				return 0, fmt.Errorf("bad format of cpu stats")
			}

			var totalClockTicks uint64
			for _, i := range fields[1:cpuFields] {
				v, err := strconvfunc.ParseUint(i)
				if err != nil {
					return 0, err
				}

				totalClockTicks += v
			}

			// 将系统累计时间以纳秒返回
			return (totalClockTicks * uint64(time.Second)) / cpuTicks, nil
		}
	}

	return 0, errors.New("bad stats format")
}
