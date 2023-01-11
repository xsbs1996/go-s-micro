package rollingwindow

// Bucket 桶
type Bucket struct {
	Sum   float64 // 数量
	Count int64   // add次数
}

// Add 桶内数据增加
func (b *Bucket) add(v float64) {
	b.Sum += v
	b.Count++
}

// Reset 桶重置
func (b *Bucket) reset() {
	b.Sum = 0
	b.Count = 0
}
