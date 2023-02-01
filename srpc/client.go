package srpc

import (
	"context"
	"github.com/sirupsen/logrus"
	discov2 "github.com/xsbs1996/go-s-micro/core/discov"
	"github.com/xsbs1996/go-s-micro/srpc/clientinterceptor"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func MustNewClient(c RpcClientConf, options ...grpc.DialOption) *grpc.ClientConn {
	c.verify()

	discov2.InitEtcdCli(clientv3.Config{Endpoints: c.Etcd.Hosts, DialTimeout: discov2.DialTimeout})
	//注册grpc解析器
	discov2.EtcdResolverRegister(discov2.NewResolver(c.Etcd.Key))

	//获取grpc连接
	cli, err := NewClient(context.Background(), c, options...)
	if err != nil {
		logrus.WithField("err", err).Fatal("MustNewClient error")
		return nil
	}

	return cli
}

func NewClient(ctx context.Context, c RpcClientConf, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	target := discov2.BuildDiscovTarget(c.Etcd.Hosts, c.Etcd.Key)
	options := buildDialOptions()
	opts = append(opts, options...)
	conn, err := grpc.DialContext(ctx, target, opts...)
	return conn, err
}

// 绑定方法
func buildDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		WithUnaryClientInterceptors(
			clientinterceptor.TracingInterceptor,
			clientinterceptor.BreakerInterceptor,
		),
	}
}

func WithUnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}
