package rollingwindow

import (
	"github.com/xsbs1996/go-s-micro/utils/timefunc"
	"sync"
	"time"
)

// RollingWindow 滑动窗口结构体
type RollingWindow struct {
	lock          sync.RWMutex  // 读写锁
	size          int           // 桶集合容量
	win           *window       // 窗口
	interval      time.Duration // 间隔时间(时间窗口大小)
	offset        int           // 当前偏移量
	lastTime      time.Duration // 最后操作时间
	ignoreCurrent bool          // 是否忽略当前桶

}

type Option func(rollingWindow *RollingWindow) // 滑动窗口操作函数

func NewRollingWindow(size int, interval time.Duration, opts ...Option) *RollingWindow {
	if size < 1 {
		panic("size must be greater than 0")
	}
	rw := &RollingWindow{
		size:     size,
		win:      newWindow(size),
		interval: interval,
		lastTime: timefunc.Now(),
	}
	for _, opt := range opts {
		opt(rw)
	}

	return rw
}

// Add 滑动窗口内桶的数据增加
func (rw *RollingWindow) Add(v float64) {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	rw.updateOffset()
	rw.win.add(rw.offset, v)
}

// Reduce 对滑动窗口内的有效桶进行操作
func (rw *RollingWindow) Reduce(fn func(b *Bucket)) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	var diff int //有效桶数量

	//获取跨度
	span := rw.span()

	//计算有效桶数量
	if span == 0 && rw.ignoreCurrent { // 如果跨越量为0且忽略当前桶,则有效桶数量为rw.size-1,这个1就是当前桶
		diff = rw.size - 1
	} else { // 如果跨越量>0,则有效桶数量为rw.size-span
		diff = rw.size - span
	}

	// 如果diff<=0,则没有任何有效桶,不操作即可
	if diff > 0 {
		offset := (rw.offset + span + 1) % rw.size // 计算当前偏移量
		rw.win.reduce(offset, diff, fn)            // offset-当前偏移量 diff-要操作的有效桶数量 timex-操作函数
	}

}

// 确定当前时间与rw.lastTime跨越了多少个桶
// 例如 lastTime = 1s, 当前时间1777ms。interval为250ms，那么跨度为3个桶
func (rw *RollingWindow) span() int {
	span := int(timefunc.Since(rw.lastTime) / rw.interval) // 跨越桶数量,向下取整

	//如果跨越桶数量在 [0,rw.size)之间则返回span
	if 0 <= span && span < rw.size {
		return span
	}

	//如果跨越桶数量<=rw.size,则返回rw.size
	return rw.size
}

func (rw *RollingWindow) updateOffset() {
	span := rw.span() // 跨越量
	if span <= 0 {    // 防止错误
		return
	}

	offset := rw.offset // 当前偏移量

	//重置已经超时的桶
	for i := 0; i < span; i++ {
		rw.win.resetBucket((offset + i + 1) % rw.size)
	}

	//更新偏移量 (老偏移量+跨越量)%size
	rw.offset = (offset + span) % rw.size

	//当前时间戳
	now := timefunc.Now()

	//减去多余时间,使rw.lastTime % rw.interval = 0,主要是为了计算span这个边界
	rw.lastTime = now - (now-rw.lastTime)%rw.interval
}
