package main

import (
	"github.com/rs/zerolog"
	"github.com/yogenyslav/ldt-2024/api/config"
	"github.com/yogenyslav/ldt-2024/api/internal/_server"
)

func main() {
	cfg := config.MustNew("./config/config.yaml")
	level, err := zerolog.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
	srv := server.New(cfg)
	srv.Run()
}
