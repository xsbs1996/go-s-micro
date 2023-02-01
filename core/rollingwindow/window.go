package rollingwindow

// Window 窗口结构体
type window struct {
	buckets []*Bucket // 桶集合
	size    int       // 桶集合容量
}

// NewWindow New一个窗口
func newWindow(size int) *window {
	buckets := make([]*Bucket, size) //创建桶集合
	//为桶集合内的桶分配内存
	for i := 0; i < size; i++ {
		buckets[i] = new(Bucket)
	}

	//return
	return &window{
		buckets: buckets,
		size:    size,
	}
}

// Add 根据偏移量增加桶内数据
func (w *window) add(offset int, v float64) {
	w.buckets[offset%w.size].add(v)
}

// Reduce 根据start起始位置使用fn操作count个桶
func (w *window) reduce(start, count int, fn func(b *Bucket)) {
	for i := 0; i < count; i++ {
		fn(w.buckets[(start+i)%w.size])
	}
}

// ResetBucket 根据偏移量重置桶
func (w *window) resetBucket(offset int) {
	w.buckets[offset%w.size].reset()
}
