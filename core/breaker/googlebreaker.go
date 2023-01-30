package breaker

import (
	"github.com/xsbs1996/go-s-micro/core/rollingwindow"
	"math"
	"sync"
	"time"
)

const (
	window     = time.Second * 10 // 10s窗口
	buckets    = 40               // 10s窗口分成40份，一份250ms
	k          = 1.5              // sre算法常数
	protection = 5                // sre算法毛刺（注:sre公式中并没有毛刺，这是结合程序实际情况添加的，防止在程序不稳定时候，如刚启动时期的错误请求会被计算在内）
	threshold  = 0.99             //熔断阈值
)

type googleBreaker struct {
	k     float64                      // sre算法常数
	stat  *rollingwindow.RollingWindow // 滑动时间窗口
	proba *proba                       // 熔断阈值
}

func newGoogleBreaker() *googleBreaker {
	interval := time.Duration(int64(window) / int64(buckets))
	rw := rollingwindow.NewRollingWindow(buckets, interval)
	return &googleBreaker{
		k:     k,
		stat:  rw,
		proba: newProba(),
	}
}

// 计算是否熔断
func (b *googleBreaker) accept() error {
	//计算滑动时间窗口内的总请求数和成功请求数量
	accepts, total := b.history()
	weightedAccepts := b.k * float64(accepts)

	//计算丢弃成功概率 sre算法
	dropRatio := math.Max(0, (float64(total-protection)-weightedAccepts)/float64(total+1))
	if dropRatio <= 0 {
		return nil
	}

	// 如果超过阈值，则触发熔断
	if b.proba.trueOnProba(dropRatio) {
		return ErrServiceUnavailable
	}
	return nil
}

// allow 手动调用
func (b *googleBreaker) allow() (internalPromise, error) {
	if err := b.accept(); err != nil {
		return nil, err
	}
	return googlePromise{
		b: b,
	}, nil
}

// doReq 自动调用
func (b *googleBreaker) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	if err := b.accept(); err != nil {
		if fallback != nil {
			return fallback(err)
		} else {
			return err
		}
	}

	defer func() {
		if e := recover(); e != nil {
			b.markFailure()
			panic(e)
		}
	}()

	err := req()
	if acceptable(err) {
		b.markSuccess()
	} else {
		b.markFailure()
	}

	return err
}

// markSuccess 调用成功则滑动窗口sum+1
func (b *googleBreaker) markSuccess() {
	b.stat.Add(1)
}

// markFailure 调用失败则滑动窗口sum+0
func (b *googleBreaker) markFailure() {
	b.stat.Add(0)
}

// 返回滑动时间窗口内的总数与成功数量
func (b *googleBreaker) history() (accepts int64, total int64) {
	b.stat.Reduce(func(b *rollingwindow.Bucket) {
		accepts += int64(b.Sum)
		total += b.Count
	})
	return
}

type googlePromise struct {
	b *googleBreaker
}

// Accept 请求成功调用
func (p googlePromise) Accept() {
	p.b.markSuccess()
}

// Reject 请求失败调用
func (p googlePromise) Reject() {
	p.b.markFailure()
}

type proba struct {
	threshold float64
	lock      sync.Mutex
}

func newProba() *proba {
	return &proba{
		threshold: threshold,
	}
}

func (p *proba) trueOnProba(proba float64) (truth bool) {
	p.lock.Lock()
	truth = p.threshold < proba
	p.lock.Unlock()
	return
}
