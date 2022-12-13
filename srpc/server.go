package srpc

import (
	"github.com/xsbs1996/go-s-micro/discov"
	"github.com/xsbs1996/go-s-micro/srpc/serverinterceptors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net"
)

// MustNewServer 综合启动项
func MustNewServer(c RpcServerConf, register GrpcRegisterFn) *RpcServer {
	c.verify()

	var options = make([]grpc.ServerOption, 0)
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		serverinterceptors.UnaryTracingInterceptor(c.Name),
	}
	options = append(options, WithUnaryServerInterceptors(unaryInterceptors...))

	server := grpc.NewServer(options...)

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
func (s *RpcServer) Start() error {
	//监听端口
	lis, err := net.Listen("tcp", s.listenOn)
	if err != nil {
		return err
	}

	//服务注册
	_, err = s.etcdRegister.Register()
	if err != nil {
		return err
	}

	//绑定grpc路由
	s.grpcRegister(s.grpcServer)

	//监听服务
	err = s.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Stop 关闭grpc服务端
func (s *RpcServer) Stop() {
	s.grpcServer.Stop()
	s.etcdRegister.Stop()
}

func WithUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}
