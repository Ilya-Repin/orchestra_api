package main

import (
	"context"
	"errors"
	app "github.com/Ilya-Repin/orchestra_api/internal/app"
	"github.com/Ilya-Repin/orchestra_api/internal/config"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/metrics"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/storage/postgres"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting orchestra api server",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	db, err := postgres.InitDB(&cfg.StorageConfig)
	if err != nil {
		panic(err)
	}

	appMetrics := metrics.New()

	application := app.NewApp(log, db, appMetrics)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPServerConfig.Port,
		ReadTimeout:  cfg.HTTPServerConfig.Timeout,
		WriteTimeout: cfg.HTTPServerConfig.Timeout,
		IdleTimeout:  cfg.HTTPServerConfig.IdleTimeout,
		Handler:      application.Routes(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("server shutdown failed", slog.Any("error", err))
		}

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("server failed", slog.Any("error", err))
	}

	<-idleConnsClosed
	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
