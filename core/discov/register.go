package discov

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"time"
)

var DialTimeout = time.Second * 3
var DefaultServiceTTL int64 = 10

type Register struct {
	ServiceName string                                  // 服务名
	ServiceAddr string                                  // 服务地址
	ServiceTTL  int64                                   // 租约时长(秒)
	leasesID    clientv3.LeaseID                        // 租约ID
	closeCh     chan struct{}                           // 关闭信号
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse // 保活通知
}

func NewRegister(ServiceName string, serviceAddr string, ServiceTTL int64) *Register {
	return &Register{
		ServiceName: ServiceName,
		ServiceAddr: serviceAddr,
		ServiceTTL:  ServiceTTL,
	}
}

// Register 服务注册
func (r *Register) Register() (chan<- struct{}, error) {
	if r.ServiceName == "" {
		return nil, errors.New("etcd service name not set")
	}
	if r.ServiceTTL <= 0 {
		return nil, errors.New("etcd service ttl setting error")
	}

	// 写入etcd
	if err := r.register(); err != nil {
		return nil, err
	}

	// 保活
	go func() {
		err := r.keepAlive()
		if err != nil {
			logrus.WithField("leasesID", r.leasesID).WithField("err", err).Error("etcd keepAlive failed")
		}
	}()

	// 设置取消信号
	r.closeCh = make(chan struct{})

	return r.closeCh, nil

}

// register 注册
func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()

	// 获取链接
	cli, err := getEtcdCli()
	if err != nil {
		return err
	}

	// 设置租约
	lease, err := cli.Grant(ctx, r.ServiceTTL)
	if err != nil {
		return err
	}
	r.leasesID = lease.ID

	// 写入etcd
	_, err = cli.Put(ctx, buildServiceName(r.ServiceName, r.leasesID), figureOutListenOn(r.ServiceAddr), clientv3.WithLease(r.leasesID))
	return err
}

// unregister 注销
func (r *Register) unregister() error {
	cli, err := getEtcdCli()
	if err != nil {
		return err
	}

	if _, err = cli.Delete(context.Background(), buildServiceName(r.ServiceName, r.leasesID)); err != nil {
		return err
	}

	if _, err = cli.Revoke(context.Background(), r.leasesID); err != nil {
		return err
	}
	return nil
}

// keepAlive 保活
func (r *Register) keepAlive() error {
	// 获取链接
	cli, err := getEtcdCli()
	if err != nil {
		return err
	}

	r.keepAliveCh, err = cli.KeepAlive(context.Background(), r.leasesID)
	if err != nil {
		return err
	}

	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				return err
			}
			return nil

		case _, ok := <-r.keepAliveCh:
			if !ok {
				if err := r.register(); err != nil {
					return err
				}
				return nil
			}
		}
	}
}

// Stop 停止保活
func (r *Register) Stop() {
	if r.closeCh != nil {
		r.closeCh <- struct{}{}
	}
}
