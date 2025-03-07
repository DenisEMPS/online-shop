package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DenisEMPS/online-shop/internal/config"
	"github.com/DenisEMPS/online-shop/internal/handler"
	"github.com/DenisEMPS/online-shop/internal/infastructure/cache"
	"github.com/DenisEMPS/online-shop/internal/infastructure/repository"
	"github.com/DenisEMPS/online-shop/internal/service"
	"github.com/DenisEMPS/online-shop/server"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	log.Info("starting application")

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		log.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	redis, err := cache.NewRedis(cfg)
	if err != nil {
		log.Error("failed to connect to redis", slog.Any("error", err))
		os.Exit(1)
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos, redis)
	handler := handler.NewHandler(service)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(cfg, handler.InitRoutes()); err != nil {
			log.Error("failed to run server", slog.Any("error", err.Error()))
			os.Exit(1)
		}
	}()

	log.Info("App started", slog.Any("port", cfg.Server.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error("error ocured on server shutting down", slog.Any("error", err))
	}

	if err := db.Close(); err != nil {
		log.Error("error ocured on db connection close", slog.Any("error", err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
