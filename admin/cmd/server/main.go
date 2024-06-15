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
// @host admin.misis.larek.tech
// @BasePath /chat
func main() {
	cfg := config.MustNew("./config.yaml")
	srv := server.New(cfg)
	srv.Run()
}
