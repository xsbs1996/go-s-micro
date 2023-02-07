package discov

import "google.golang.org/grpc/resolver"

// Scheme 返回此解析器支持的方案 用于grpc
func (r *Resolver) Scheme() string {
	return r.schema
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

// ResolveNow resolver.Resolver interface 用于grpc
func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {}

// Close resolver.Resolver 取消 用于grpc
func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}
