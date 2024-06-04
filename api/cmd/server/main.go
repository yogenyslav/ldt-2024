package main

import (
	"github.com/yogenyslav/ldt-2024/api/config"
	"github.com/yogenyslav/ldt-2024/api/internal/_server"
)

func main() {
	cfg := config.MustNew("./config/config.yaml")
	srv := server.New(cfg)
	srv.Run()
}
