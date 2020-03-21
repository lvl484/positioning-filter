package consul

import (
	"fmt"
	"net/http"

	consulApi "github.com/hashicorp/consul/api"
)

const (
	agentDeregisterAfter = 45
	agentCheckInterval   = 15
	agentCheckTimeout    = 5
)

// AgentConfig returns new consul agent config
func (c *Config) AgentConfig() *consulApi.AgentServiceRegistration {
	return &consulApi.AgentServiceRegistration{
		ID:   c.ServiceName,
		Name: c.ServiceName,
		Check: &consulApi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d%s", c.ServiceName, c.ServicePort, c.ServiceHealthCheckPath),
			Method:                         http.MethodGet,
			Interval:                       fmt.Sprintf("%ds", agentCheckInterval),
			Timeout:                        fmt.Sprintf("%ds", agentCheckTimeout),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", agentDeregisterAfter),
		},
		Address: c.ServiceName,
		Port:    c.ServicePort,
	}
}
