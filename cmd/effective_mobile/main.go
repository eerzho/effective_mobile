package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "effective_mobile/docs"
	"effective_mobile/internal/app"
	"effective_mobile/internal/config"
	"effective_mobile/internal/conn"
	"effective_mobile/internal/lib/logger/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Effective mobile
// @version 1.0
// @description test task for effective mobile
// @contact.name Zhanbolat
// @contact.email eerzho@gmail.com
// @basePath /api
func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Debug("debug message enabled")

	log.Info("starting connection to postgres")

	connection := conn.MustNew(cfg)

	log.Info("connected to postgres", slog.String("url", cfg.Postgres.Url))

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Error("failed to close postgres connection", sl.Err(err))

			return
		}
		log.Info("closed connection")
	}(connection.Psql.DB)

	log.Info("starting http server", slog.String("address", cfg.Address))

	application := app.New(log, connection)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	httpSrv := &http.Server{
		Addr:         cfg.Address,
		Handler:      application.HttpApp.Router,
		ReadTimeout:  cfg.Http.Timeout,
		WriteTimeout: cfg.Http.Timeout,
		IdleTimeout:  cfg.Http.IdleTimeout,
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Error("failed to start http server", sl.Err(err))
		}
	}()

	log.Info("started http server")

	<-done
	log.Info("stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("stopped http server")
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
	}

	return log
}
