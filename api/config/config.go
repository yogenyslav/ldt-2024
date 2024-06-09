package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	srvconfig "github.com/yogenyslav/ldt-2024/api/internal/_server/config"
	"github.com/yogenyslav/ldt-2024/api/pkg/client"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
)

// Config holds the configuration for the application.
type Config struct {
	Server     *srvconfig.Config        `yaml:"server"`
	Postgres   *postgres.Config         `yaml:"postgres"`
	Jaeger     *tracing.Config          `yaml:"jaeger"`
	Prometheus *prom.Config             `yaml:"prometheus"`
	KeyCloak   *KeyCloakConfig          `yaml:"keycloak"`
	Prompter   *client.GrpcClientConfig `yaml:"prompter"`
}

// KeyCloakConfig holds the configuration for Keycloak.
type KeyCloakConfig struct {
	URL          string `yaml:"url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Realm        string `yaml:"realm"`
	AdminRealm   string `yaml:"admin_realm"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
}

// MustNew reads the configuration from the given path and returns a Config struct
// or panics if failed.
func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Panic().Err(err).Msg("failed to read config")
	}
	return cfg
}
