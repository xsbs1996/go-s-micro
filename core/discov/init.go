package discov

import (
	"errors"
	"fmt"
	_ "github.com/xsbs1996/go-s-micro/logsj"
	"go.etcd.io/etcd/client/v3"
	"sync"
)

var (
	etcdCli *clientv3.Client
	once    sync.Once
)

// InitEtcdCli 初始化etcd链接
func InitEtcdCli(config clientv3.Config) {
	var err error
	once.Do(func() {
		etcdCli, err = clientv3.New(config)
		if err != nil {
			panic(fmt.Sprintf("%s:%v", "Failed to start etcd", err))
		}
	})
}

// 获取etcd链接
func getEtcdCli() (*clientv3.Client, error) {
	if etcdCli == nil {
		return nil, errors.New("etcd not initialized")
	}

	return etcdCli, nil
}
