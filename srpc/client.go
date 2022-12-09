package srpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xsbs1996/go-s-micro/discov"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

func MustNewClient(c RpcClientConf, options ...grpc.DialOption) *grpc.ClientConn {
	c.verify()

	discov.InitEtcdCli(clientv3.Config{Endpoints: c.Etcd.Hosts, DialTimeout: discov.DialTimeout})
	//注册grpc解析器
	discovResolver := discov.NewResolver(c.Etcd.Key)
	resolver.Register(discovResolver)

	//获取grpc连接
	cli, err := NewClient(context.Background(), c, options...)
	if err != nil {
		logrus.WithField("err", err).Error("MustNewClient error")
	}

	return cli
}

func NewClient(ctx context.Context, c RpcClientConf, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	target := discov.BuildDiscovTarget(c.Etcd.Hosts, c.Etcd.Key)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.DialContext(ctx, target, opts...)
	return conn, err
}
