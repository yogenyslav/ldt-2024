package main

import (
	"github.com/rs/zerolog"
	"github.com/yogenyslav/ldt-2024/api/config"
	"github.com/yogenyslav/ldt-2024/api/internal/_server"
)

// @title Chat service API
// @version 1.0
// @description Документация API чат-сервиса команды misis.tech
// @license.name BSD-3-Clause
// @license.url https://opensource.org/license/bsd-3-clause
// @host localhost:10000
// @BasePath /
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
