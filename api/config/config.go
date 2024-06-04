package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	srvconfig "github.com/yogenyslav/ldt-2024/api/internal/_server/config"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
)

type Config struct {
	Server     *srvconfig.Config `yaml:"server"`
	Postgres   *postgres.Config  `yaml:"postgres"`
	Jaeger     *tracing.Config   `yaml:"jaeger"`
	Prometheus *prom.Config      `yaml:"prometheus"`
	KeyCloak   *KeyCloakConfig   `yaml:"keycloak"`
}

type KeyCloakConfig struct {
	KeyCloakURL          string `yaml:"keycloak_url"`
	KeyCloakClientID     string `yaml:"keycloak_client_id"`
	KeyCloakClientSecret string `yaml:"keycloak_client_secret"`
	KeyCloakRealm        string `yaml:"keycloak_realm"`
}

func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Panic().Err(err).Msg("failed to read config")
	}
	return cfg
}
