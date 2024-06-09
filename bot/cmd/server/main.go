package main

import (
	"github.com/rs/zerolog"
	"github.com/yogenyslav/ldt-2024/bot/config"
	server "github.com/yogenyslav/ldt-2024/bot/internal/_server"
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
	srv := server.New(cfg)
	srv.Run()
}
