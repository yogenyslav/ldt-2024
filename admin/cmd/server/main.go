package main

import (
	"github.com/yogenyslav/ldt-2024/admin/config"
	server "github.com/yogenyslav/ldt-2024/admin/internal/_server"
)

// @title Admin service API
// @version 1.0
// @description Документация API админ-сервиса команды misis.tech
// @license.name BSD-3-Clause
// @license.url https://opensource.org/license/bsd-3-clause
// @host api.misis.larek.tech
// @BasePath /admin
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.MustNew("./config/config.yaml")
	srv := server.New(cfg)
	srv.Run()
}
