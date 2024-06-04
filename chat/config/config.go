package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	srvconfig "github.com/yogenyslav/ldt-2024/chat/internal/_server/config"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
)

// Config is the struct for the whole service configuration.
type Config struct {
	Server     *srvconfig.Config `yaml:"server"`
	Postgres   *postgres.Config  `yaml:"postgres"`
	Jaeger     *tracing.Config   `yaml:"jaeger"`
	Prometheus *prom.Config      `yaml:"prometheus"`
}

// MustNew creates a new Config instance or panics if failed.
func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Panic().Err(err).Msg("failed to parse config")
	}
	return cfg
}
