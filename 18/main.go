package main

import (
	"calendar/internal/config"
	"calendar/internal/handlers"
	sl2 "calendar/internal/lib/log"
	"calendar/internal/service"
	"calendar/internal/storage"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	stg := storage.New()

	srv := service.New(stg)

	h := handlers.New(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", h.Create())
	mux.HandleFunc("/update_event", h.Update())
	mux.HandleFunc("/delete_event", h.Delete())
	mux.HandleFunc("/events_for_day", h.GetForDay())
	mux.HandleFunc("/events_for_week", h.GetForWeek())
	mux.HandleFunc("/events_for_month", h.GetForMonth())

	wrMux := h.Logger(mux)

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      wrMux,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
		}
	}()

	slog.Info("server started", slog.String("address", cfg.HTTPServer.Address))

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}
	slog.Info("shut down")

}
