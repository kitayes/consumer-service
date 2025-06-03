package redis

import (
	"consumer-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Cache struct {
	client *redis.Client
	config Config
	logger Logger
}

func NewCache(cfg Config, logger Logger) *Cache {
	return &Cache{
		config: cfg,
		logger: logger,
	}
}

func (c *Cache) Init() error {
	c.client = redis.NewClient(&redis.Options{
		Addr:     c.config.Addr,
		Password: c.config.Password,
		DB:       c.config.DB,
	})

	c.logger.Info("Redis cache initialized at %s", c.config.Addr)
	return nil
}

func (c *Cache) Run(_ context.Context) {}

func (c *Cache) Stop() {
	if err := c.client.Close(); err != nil {
		c.logger.Warn("Failed to close Redis: %v", err)
	}
}

func (c *Cache) SaveOrder(ctx context.Context, order *models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	key := fmt.Sprintf("order:%d", order.ID)
	if err := c.client.Set(ctx, key, data, time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to save order in Redis: %w", err)
	}

	c.logger.Debug("Order cached in Redis: key=%s", key)
	return nil
}
