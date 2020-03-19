package consul

import consulApi "github.com/hashicorp/consul/api"

// NewClient new consul client
func (c *Config) NewClient() (*consulApi.Client, error) {
	consulCfg := consulApi.DefaultConfig()
	consulCfg.Address = c.Address

	return consulApi.NewClient(consulCfg)
}
