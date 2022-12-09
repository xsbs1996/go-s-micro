package discov

type EtcdRegisterConf struct {
	Hosts []string
	Key   string
	User  string
	Pass  string
}

type EtcdResolverConf struct {
	Hosts []string
	Key   string
	User  string
	Pass  string
}
