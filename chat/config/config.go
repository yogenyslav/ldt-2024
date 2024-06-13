package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	srvconfig "github.com/yogenyslav/ldt-2024/chat/internal/_server/config"
	"github.com/yogenyslav/ldt-2024/chat/pkg/client"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
	rediscache "github.com/yogenyslav/pkg/storage/redis_cache"
)

// Config конфигурация сервиса.
type Config struct {
	Server   *srvconfig.Config        `yaml:"server"`
	Postgres *postgres.Config         `yaml:"postgres"`
	Jaeger   *tracing.Config          `yaml:"jaeger"`
	Redis    *rediscache.Config       `yaml:"redis"`
	ChatProm *prom.Config             `yaml:"chat_prom"`
	BotProm  *prom.Config             `yaml:"bot_prom"`
	API      *client.GrpcClientConfig `yaml:"api"`
	KeyCloak *KeyCloakConfig          `yaml:"keycloak"`
}

// KeyCloakConfig конфигурация KeyCloak.
type KeyCloakConfig struct {
	URL   string `yaml:"url"`
	Realm string `yaml:"realm"`
}

// MustNew создает новую конфигурацию или вызывает панику.
func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Panic().Err(err).Msg("failed to parse config")
	}
	return cfg
}
