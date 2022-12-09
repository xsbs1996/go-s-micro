package srpc

import (
	"github.com/xsbs1996/go-s-micro/discov"
	"google.golang.org/grpc"
)

type RpcServer struct {
	listenOn     string
	timeout      int64
	grpcServer   *grpc.Server
	grpcRegister GrpcRegisterFn
	etcdRegister *discov.Register
}

type RpcServerConf struct {
	ListenOn string
	Etcd     discov.EtcdRegisterConf
	Timeout  int64
}

type RpcClientConf struct {
	Etcd    discov.EtcdResolverConf
	Timeout int64
}

type GrpcRegisterFn func(*grpc.Server)

func (c RpcServerConf) verify() {
	if c.ListenOn == "" {
		panic("Please set ListenOn")
	}

	if c.Etcd.Hosts == nil || c.Etcd.Key == "" {
		panic("Please set etcd")
	}
}

func (c RpcClientConf) verify() {
	if c.Etcd.Hosts == nil || c.Etcd.Key == "" {
		panic("Please set etcd")
	}
}
