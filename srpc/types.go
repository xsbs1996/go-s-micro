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
	Name     string                  `yaml:"Name" required:"true"`
	ListenOn string                  `yaml:"ListenOn" required:"true"`
	Etcd     discov.EtcdRegisterConf `yaml:"Etcd"`
	Timeout  int64                   `yaml:"Timeout" default:"2000"`
}

type RpcClientConf struct {
	Etcd    discov.EtcdResolverConf `yaml:"Etcd" required:"true"`
	Timeout int64                   `yaml:"Timeout" default:"2000"`
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
