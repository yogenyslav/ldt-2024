package main

import (
	"github.com/yogenyslav/ldt-2024/chat/config"
	"github.com/yogenyslav/ldt-2024/chat/internal/_server"
)

func main() {
	cfg := config.MustNew("./config/config.yaml")
	srv := server.New(cfg)
	srv.Run()
}
