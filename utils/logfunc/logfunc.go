package logfunc

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const CallerDepth = 4

const timeFormat = "2006-01-02T15:04:05.000Z07:00"

func GetCaller(callDepth int) string {
	_, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		return ""
	}

	return prettyCaller(file, line)
}

func GetTimestamp() string {
	return time.Now().Format(timeFormat)
}

func prettyCaller(file string, line int) string {
	idx := strings.LastIndexByte(file, '/')
	if idx < 0 {
		return fmt.Sprintf("%s:%d", file, line)
	}

	idx = strings.LastIndexByte(file[:idx], '/')
	if idx < 0 {
		return fmt.Sprintf("%s:%d", file, line)
	}

	return fmt.Sprintf("%s:%d", file[idx+1:], line)
}
