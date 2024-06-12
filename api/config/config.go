package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	srvconfig "github.com/yogenyslav/ldt-2024/api/internal/_server/config"
	"github.com/yogenyslav/ldt-2024/api/pkg/client"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/mongo"
	"github.com/yogenyslav/pkg/storage/postgres"
)

// Config конфигурация сервиса.
type Config struct {
	Server     *srvconfig.Config        `yaml:"server"`
	Postgres   *postgres.Config         `yaml:"postgres"`
	Jaeger     *tracing.Config          `yaml:"jaeger"`
	Prometheus *prom.Config             `yaml:"prometheus"`
	KeyCloak   *KeyCloakConfig          `yaml:"keycloak"`
	Prompter   *client.GrpcClientConfig `yaml:"prompter"`
	Predictor  *client.GrpcClientConfig `yaml:"predictor"`
	Mongo      *mongo.Config            `yaml:"mongo"`
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

// MustNew читает конфигурацию из файла и возвращает ее или panic.
func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Panic().Err(err).Msg("failed to read config")
	}
	return cfg
}
