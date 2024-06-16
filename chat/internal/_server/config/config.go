package srvconfig

// Config конфигурация сервера.
type Config struct {
	CipherKey   string   `yaml:"cipher_key"`
	LogLevel    string   `yaml:"log_level"`
	CorsOrigins []string `yaml:"cors_origins"`
	Port        int      `yaml:"port"`
	BotToken    string   `yaml:"bot_token"`
	DebugMode   bool     `yaml:"debug_mode"`
	AppURL      string   `yaml:"app_url"`
}
