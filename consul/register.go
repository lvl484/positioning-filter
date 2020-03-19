package consul

import consulApi "github.com/hashicorp/consul/api"

// ServiceRegister is used to register a new service with the agent
func (c *Config) ServiceRegister(client *consulApi.Client, agent *consulApi.AgentServiceRegistration) error {
	return client.Agent().ServiceRegister(agent)
}
