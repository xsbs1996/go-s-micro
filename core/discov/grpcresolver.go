package discov

import "google.golang.org/grpc/resolver"

// EtcdResolverRegister 注册grpc解析器
func EtcdResolverRegister(r *Resolver) {
	resolver.Register(r)
}
