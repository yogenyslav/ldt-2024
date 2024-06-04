package srvconfig

// Config is the struct for the server configuration.
type Config struct {
	LogLevel    string   `yaml:"log_level"`
	CorsOrigins []string `yaml:"cors_origins"`
	Port        int      `yaml:"port"`
}
