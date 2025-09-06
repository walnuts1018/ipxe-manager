package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Code-Hex/synchro/tz"
	appConfig "github.com/walnuts1018/ipxe-manager/config"
	"github.com/walnuts1018/ipxe-manager/infrastructure/database"
	"github.com/walnuts1018/ipxe-manager/logger"
	"github.com/walnuts1018/ipxe-manager/tracer"
	"github.com/walnuts1018/ipxe-manager/util/clock"
	"github.com/walnuts1018/ipxe-manager/wire"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	cfg, err := appConfig.Load()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger := logger.CreateLogger(cfg.LogLevel, cfg.LogType)
	slog.SetDefault(logger)

	closeTracer, err := tracer.NewTracerProvider(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create tracer provider", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer closeTracer()

	db, err := database.NewBoltDB(cfg.DB)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close database", slog.String("error", err.Error()))
		}
	}()

	router, err := wire.CreateRouter(ctx, cfg, db, clock.Default[tz.AsiaTokyo]())
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create router", slog.String("error", err.Error()))
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: router,
	}

	go func() {
		slog.Info("Server is running", slog.String("port", cfg.App.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "Failed to run server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	stop()
	slog.Info("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
