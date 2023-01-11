package timefunc

import "time"

// 设定初始时间,这种方式类似于设定初始的时间戳时间
var initTime = GetThisYearZero().AddDate(-1, 0, 0)

// Now 当前时间-initTime
func Now() time.Duration {
	return time.Since(initTime)
}

// Since 计算当前时间戳与给定时间戳d的时间跨度
func Since(d time.Duration) time.Duration {
	now := time.Since(initTime) // 当前时间戳
	return now - d              // 时间跨度
}

func Time() time.Time {
	return initTime.Add(Now())
}

// GetThisYearZero 获取今年零点
func GetThisYearZero() time.Time {
	timeNow := time.Now()
	return time.Date(timeNow.Year(), 1, 1, 0, 0, 0, 0, timeNow.Location())
}
