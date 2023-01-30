package discov

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

const schema = "etcd"

// Resolver 服务发现
type Resolver struct {
	schema         string              // 解析器方案
	keyPrefix      string              // key前缀
	serviceName    string              // 服务名
	serverList     map[string]string   // 服务列表
	watchCh        clientv3.WatchChan  // 监控服务
	onceLock       *sync.Mutex         // 互斥锁
	once           bool                // 单例
	lock           *sync.RWMutex       // 读写锁
	closeCh        chan struct{}       // 关闭信号
	grpcClientConn resolver.ClientConn // grpc resolver.ClientConn
	grpcAddrsList  []resolver.Address  // grpc resolver.Address
}

// NewResolver  新建发现服务
func NewResolver(ServiceName string) *Resolver {
	return &Resolver{
		schema:      schema,
		serviceName: ServiceName,
		serverList:  make(map[string]string, 0),
		lock:        new(sync.RWMutex),
		onceLock:    new(sync.Mutex),
		once:        false,
	}
}

// Start 初始化服务列表和监视
func (r *Resolver) Start() (chan<- struct{}, error) {
	r.onceLock.Lock()
	defer r.onceLock.Unlock()
	if r.once {
		return r.closeCh, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()

	if r.serviceName == "" {
		return nil, errors.New("etcd service name not set")
	}

	if r.keyPrefix == "" {
		r.keyPrefix = r.serviceName
	}

	// 获取链接
	cli, err := getEtcdCli()
	if err != nil {
		return nil, err
	}

	resp, err := cli.Get(ctx, r.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	//遍历获取到的key和value
	for _, ev := range resp.Kvs {
		r.setServiceList(string(ev.Key), string(ev.Value))
	}

	//监视前缀，修改变更的server
	go func() {
		err := r.watcher()
		if err != nil {
			logrus.WithField("err", err).Error("etcd watcher failed")
		}
	}()

	// 设置取消信号
	r.closeCh = make(chan struct{})
	r.once = true

	return r.closeCh, nil
}

// watcher 监听key的前缀变化
func (r *Resolver) watcher() error {
	// 获取链接
	cli, err := getEtcdCli()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()

	r.watchCh = cli.Watch(context.Background(), r.keyPrefix, clientv3.WithPrefix())
	for {
		select {
		case <-r.closeCh:
			return nil
		case res, ok := <-r.watchCh:
			if ok {
				r.updateServiceList(res.Events)
			} else {
				if _, err = r.Start(); err != nil {
					return err
				}
			}
		case <-ticker.C:
			r.sync()
		}
	}
}

// sync 同步获取所有地址信息
func (r *Resolver) sync() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 获取链接
	cli, err := getEtcdCli()
	if err != nil {
		return
	}

	res, err := cli.Get(ctx, r.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return
	}

	for _, v := range res.Kvs {
		r.setServiceList(parseValue(v.Key), parseValue(v.Value))
	}

	if r.grpcClientConn != nil {
		if err := r.grpcClientConn.UpdateState(resolver.State{Addresses: r.grpcAddrsList}); err != nil {
			logrus.WithField("resolver.State", r.grpcAddrsList).WithField("err", err).Fatal("etcd sync UpdateState failed")
			return
		}
	}

	return
}

// updateServiceList 更新grpc服务端地址列表
func (r *Resolver) updateServiceList(events []*clientv3.Event) {
	for _, ev := range events {
		switch ev.Type {
		case clientv3.EventTypePut:
			r.setServiceList(parseValue(ev.Kv.Key), parseValue(ev.Kv.Value))
		case clientv3.EventTypeDelete:
			r.delServiceList(parseValue(ev.Kv.Key))
		}

		if r.grpcClientConn != nil {
			if err := r.grpcClientConn.UpdateState(resolver.State{Addresses: r.grpcAddrsList}); err != nil {
				logrus.WithField("key", ev.Kv.Key).WithField("value", ev.Kv.Value).WithField("err", err).Fatal("etcd UpdateState failed")
				return
			}
		}
	}

}

// SetServiceList 新增服务地址
func (r *Resolver) setServiceList(key, val string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.serverList[key] = val
	addr := resolver.Address{Addr: val}
	if !existAddr(r.grpcAddrsList, addr) {
		r.grpcAddrsList = append(r.grpcAddrsList, addr)
	}

}

// DelServiceList 删除服务地址
func (r *Resolver) delServiceList(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	value, ok := r.serverList[key]
	if !ok {
		return
	}
	delete(r.serverList, key)

	for index, addr := range r.grpcAddrsList {
		if addr.Addr == value {
			r.grpcAddrsList = append(r.grpcAddrsList[:index], r.grpcAddrsList[index+1:]...)
		}
	}

}

// GetServices 获取服务地址
func (r *Resolver) GetServices() map[string]string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.serverList
}

// GetKeyPrefix 返回前缀
func (r *Resolver) GetKeyPrefix() string {
	return r.keyPrefix
}

// EtcdResolverRegister 注册grpc解析器
func EtcdResolverRegister(r *Resolver) {
	resolver.Register(r)
}

// Stop 关闭服务
func (r *Resolver) Stop() {
	r.closeCh <- struct{}{}
}

// Scheme 返回此解析器支持的方案 用于grpc
func (r *Resolver) Scheme() string {
	return r.schema
}

// ResolveNow resolver.Resolver interface 用于grpc
func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {}

// Close resolver.Resolver 取消 用于grpc
func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

// Build 为给定的目标创建一个新的 resolver.Resolver 用于grpc
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.grpcClientConn = cc

	r.keyPrefix = getEndpoints(target)
	r.serviceName = r.keyPrefix

	if _, err := r.Start(); err != nil {
		return nil, err
	}
	return r, nil
}
