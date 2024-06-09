package srvconfig

// Config is the struct for the server configuration.
type Config struct {
	LogLevel    string `yaml:"log_level"`
	Port        int    `yaml:"port"`
	GatewayPort int    `yaml:"gateway_port"`
}
