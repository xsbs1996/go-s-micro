package cgroup

import (
	"fmt"
	"github.com/xsbs1996/go-s-micro/core/iosj"
	"github.com/xsbs1996/go-s-micro/core/logsj"
	"github.com/xsbs1996/go-s-micro/utils/strconvfunc"
	"os"
	"path"
	"strings"
)

const (
	cgroupDir = "/sys/fs/cgroup" // cgroup目录
)
const (
	cgroupCpuacctUsage       = "cpuacct.usage"        // 所有task的cpu使用总时长
	cgroupCpuacctUsagePercpu = "cpuacct.usage_percpu" // 所有task使用每个cpu的时长
	cgroupCpuCfsQuotaUs      = "cpu.cfs_quota_us"     // 周期内最多可使用的时间。 注：若cfs_quota_us是cfs_period_us的两倍，就表示在两个核上完全使用，数值范围为 1000 - 1000,000（微秒）
	cgroupCpuCfsPeriodUs     = "cpu.cfs_period_us"    // 周期时间。 注：必须与cfs_quota_us配合使用
	cgroupCpusetCpus         = "cpuset.cpus"          // cgroup可使用的CPU编号，如0-2,16代表 0、1、2 和 16 这 4 个 CPU。
)

var cgroupFile = fmt.Sprintf("/proc/%d/cgroup", os.Getpid()) // cgroup子系统对应名称与目录

var cgroups = make(map[string]string, 0) // cgroup子系统对应名称与目录

// 获取cgroup目录下资源名称存放文件
func init() {
	err := getResourceList()
	if err != nil {
		logsj.CgroupLog(err)
		return
	}
}

// 获取cgroup子系统对应名称与目录
func getResourceList() error {
	lines, err := iosj.ReadTextLines(cgroupFile, iosj.WithoutBlank())
	if err != nil {
		return err
	}

	for _, line := range lines {
		cols := strings.Split(line, ":")
		if len(cols) != 3 {
			return fmt.Errorf("invalid cgroup line: %s", line)
		}

		subSys := cols[1]
		if !strings.HasPrefix(subSys, "cpu") {
			continue
		}

		cgroups[subSys] = path.Join(cgroupDir, subSys)
		if strings.Contains(subSys, ",") {
			for _, k := range strings.Split(subSys, ",") {
				cgroups[k] = path.Join(cgroupDir, k)
			}
		}
	}

	return nil
}

// GetCpuacctUsage 获取所有task的cpu使用总时长
func GetCpuacctUsage() (uint64, error) {
	filePath, ok := cgroups["cpuacct"]
	if !ok {
		return 0, nil
	}
	data, err := iosj.ReadText(path.Join(filePath, cgroupCpuacctUsage))
	if err != nil {
		return 0, err
	}

	return strconvfunc.ParseUint(data)
}

// GetCpuacctUsagePercpu 获取所有task使用每个cpu的时长
func GetCpuacctUsagePercpu() ([]uint64, error) {
	filePath, ok := cgroups["cpuacct"]
	if !ok {
		return nil, nil
	}

	data, err := iosj.ReadText(path.Join(filePath, cgroupCpuacctUsagePercpu))
	if err != nil {
		return nil, err
	}

	var usage []uint64
	for _, v := range strings.Fields(string(data)) {
		u, err := strconvfunc.ParseUint(v)
		if err != nil {
			return nil, err
		}

		usage = append(usage, u)
	}

	return usage, nil

}

// GetCpuCfsQuotaUs 获取周期内最多可使用的时间
func GetCpuCfsQuotaUs() (uint64, error) {
	filePath, ok := cgroups["cpu"]
	if !ok {
		return 0, nil
	}

	data, err := iosj.ReadText(path.Join(filePath, cgroupCpuCfsQuotaUs))
	if err != nil {
		return 0, err
	}

	return strconvfunc.ParseUint(data)
}

// GetCpuCfsPeriodUs 获取周期时间
func GetCpuCfsPeriodUs() (uint64, error) {
	filePath, ok := cgroups["cpu"]
	if !ok {
		return 0, nil
	}

	data, err := iosj.ReadText(path.Join(filePath, cgroupCpuCfsPeriodUs))
	if err != nil {
		return 0, err
	}

	return strconvfunc.ParseUint(data)
}

// GetCpusetCpus 获取可使用的CPU编号
func GetCpusetCpus() ([]uint64, error) {
	filePath, ok := cgroups["cpuset"]
	if !ok {
		return nil, nil
	}

	data, err := iosj.ReadText(path.Join(filePath, cgroupCpusetCpus))
	if err != nil {
		return nil, err
	}

	return strconvfunc.ParseUints(data)
}
