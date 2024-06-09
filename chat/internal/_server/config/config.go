package srvconfig

// Config is the struct for the server configuration.
type Config struct {
	CipherKey   string   `yaml:"cipher_key"`
	LogLevel    string   `yaml:"log_level"`
	CorsOrigins []string `yaml:"cors_origins"`
	Port        int      `yaml:"port"`
}
