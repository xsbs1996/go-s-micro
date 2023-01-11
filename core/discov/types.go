package discov

type EtcdRegisterConf struct {
	Hosts []string `yaml:"Hosts" required:"true"`
	Key   string   `yaml:"Key" required:"true"`
	User  string   `yaml:"User"`
	Pass  string   `yaml:"Pass"`
}

type EtcdResolverConf struct {
	Hosts []string `yaml:"Hosts" required:"true"`
	Key   string   `yaml:"Key" required:"true"`
	User  string   `yaml:"User"`
	Pass  string   `yaml:"Pass"`
}
