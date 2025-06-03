package main

import (
	"consumer-service/internal/application"
	"consumer-service/internal/deliviery/broker"
	"consumer-service/internal/repository/postgres"
	"consumer-service/internal/repository/redis"
	"consumer-service/pkg/config"
	"consumer-service/pkg/logger"
	service "consumer-service/pkg/services"
	"context"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log/slog"
	"os"
)

type Config struct {
	Repo   postgres.Config `envPrefix:"REPO_"`
	Logger logger.Config   `envPrefix:"LOGGER_"`
	Cache  redis.Config    `envPrefix:"REDIS_"`
	Broker broker.Config   `envPrefix:"BROKER_"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load .env!", slog.Any("error", err))
		os.Exit(1)
	}

	cfg := Config{}
	if err := config.ReadEnvConfig(&cfg); err != nil {
		slog.Error("Failed to read environment configuration!", slog.Any("error", err))
		os.Exit(1)
	}

	log := logger.NewLogger(&cfg.Logger)

	cache := redis.NewCache(cfg.Cache, log)
	repos := postgres.NewPostgresRepository(&cfg.Repo, log)
	orderService := application.NewOrderService(repos, cache)
	kafka := broker.NewConsumer(&cfg.Broker, log, orderService)

	srv := service.NewManager(log)
	srv.AddService(
		repos,
		kafka,
		cache,
		orderService,
	)

	ctx := context.Background()
	if err := srv.Run(ctx); err != nil {
		err := errors.Wrap(err, "srv.Run err:")
		log.Error(err.Error())
		return
	}

	log.Info("consumer-service started")
}
