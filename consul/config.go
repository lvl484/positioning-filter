package consul

// Config for consul
type Config struct {
	Address                string
	ServiceName            string
	ServicePort            int
	ServiceHealthCheckPath string
}
