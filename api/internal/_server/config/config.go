package srvconfig

// Config конфигурация сервера.
type Config struct {
	LogLevel    string `yaml:"log_level"`
	Port        int    `yaml:"port"`
	GatewayPort int    `yaml:"gateway_port"`
}
