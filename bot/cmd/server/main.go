package main

import (
	"github.com/rs/zerolog"
	"github.com/yogenyslav/ldt-2024/bot/config"
	"github.com/yogenyslav/pkg/loctime"
)

func main() {
	cfg := config.MustNew("./config/config.yaml")
	level, err := zerolog.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
	loctime.SetLocation(loctime.MoscowLocation)
}
