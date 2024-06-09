package srvconfig

type ServerConfig struct {
	BotToken  string `yaml:"bot_token"`
	LogLevel  string `yaml:"log_level"`
	DebugMode bool   `yaml:"debug_mode"`
	Port      int    `yaml:"port"`
}
