package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	srvconfig "github.com/yogenyslav/ldt-2024/admin/internal/_server/config"
	"github.com/yogenyslav/ldt-2024/admin/pkg/client"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/minios3"
	"github.com/yogenyslav/pkg/storage/postgres"
)

// Config конфигурация сервиса.
type Config struct {
	Server   *srvconfig.Config        `yaml:"server"`
	Postgres *postgres.Config         `yaml:"postgres"`
	S3       *minios3.Config          `yaml:"s3"`
	Jaeger   *tracing.Config          `yaml:"jaeger"`
	Prom     *prom.Config             `yaml:"prom"`
	API      *client.GrpcClientConfig `yaml:"api"`
	KeyCloak *KeyCloakConfig          `yaml:"keycloak"`
}

// KeyCloakConfig конфигурация KeyCloak.
type KeyCloakConfig struct {
	URL          string `yaml:"url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Realm        string `yaml:"realm"`
	AdminRealm   string `yaml:"admin_realm"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
}

// MustNew создает новую конфигурацию или вызывает панику.
func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Panic().Err(err).Msg("failed to parse config")
	}
	return cfg
}
