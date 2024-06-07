package srvconfig

type ServerConfig struct {
	BotToken string `yaml:"bot_token"`
	LogLevel string `yaml:"log_level"`
}
