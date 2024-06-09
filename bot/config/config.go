package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	srvconfig "github.com/yogenyslav/ldt-2024/bot/internal/_server/config"
	"github.com/yogenyslav/ldt-2024/bot/pkg/client"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
	rediscache "github.com/yogenyslav/pkg/storage/redis_cache"
)

// Config is a configuration for the bot.
type Config struct {
	Server     *srvconfig.ServerConfig  `yaml:"server"`
	Postgres   *postgres.Config         `yaml:"postgres"`
	Jaeger     *tracing.Config          `yaml:"jaeger"`
	Prometheus *prom.Config             `yaml:"prometheus"`
	API        *client.GrpcClientConfig `yaml:"api"`
	Redis      *rediscache.Config       `yaml:"redis"`
	KeyCloak   *KeyCloakConfig          `yaml:"keycloak"`
}

// KeyCloakConfig holds the configuration for Keycloak.
type KeyCloakConfig struct {
	URL   string `yaml:"url"`
	Realm string `yaml:"realm"`
}

// MustNew reads the configuration from the file and returns it or panics if failed.
func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Panic().Err(err).Msg("failed to read config")
	}
	return cfg
}
