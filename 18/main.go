package main

import (
	"calendar/internal/config"
	sl2 "calendar/internal/lib/log"
	"calendar/internal/storage"
	"log"
	"log/slog"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sl := sl2.SetupLogger(cfg.Env)
	slog.SetDefault(sl)
	sl.Info("config loaded, start application")

	// TODO: graceful

	storage := storage.New()

	// TODO: service

	// TODO handlers

	// TODO: SERVER and MiddleWare

}
