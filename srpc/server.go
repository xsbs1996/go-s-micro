package srpc

import (
	"github.com/xsbs1996/go-s-micro/discov"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net"
)

// MustNewServer 综合启动项
func MustNewServer(c RpcServerConf, register GrpcRegisterFn) *RpcServer {
	server := grpc.NewServer()
	c.verify()

	discov.InitEtcdCli(clientv3.Config{Endpoints: c.Etcd.Hosts, DialTimeout: discov.DialTimeout})
	etcdRegister := discov.NewRegister(c.Etcd.Key, c.ListenOn, discov.DefaultServiceTTL)

	return &RpcServer{
		listenOn:     c.ListenOn,
		timeout:      c.Timeout,
		grpcServer:   server,
		etcdRegister: etcdRegister,
		grpcRegister: register,
	}
}

// Start 启动grpc服务端
func (s *RpcServer) Start() {
	//监听端口
	lis, err := net.Listen("tcp", s.listenOn)
	if err != nil {
		panic(err)
	}

	//服务注册
	_, err = s.etcdRegister.Register()
	if err != nil {
		panic(err)
	}

	//绑定grpc路由
	s.grpcRegister(s.grpcServer)

	//监听服务
	err = s.grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

// Stop 关闭grpc服务端
func (s *RpcServer) Stop() {
	s.grpcServer.Stop()
	s.etcdRegister.Stop()
}
