package broker

import (
	"consumer-service/internal/models"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string `env:"BROKER_BROKERS" envSeparator:","`
	Topic   string   `env:"BROKER_TOPIC"`
}

type OrderProcessor interface {
	ProcessOrder(ctx context.Context, order *models.Order) error
}

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Consumer struct {
	reader    *kafka.Reader
	config    *Config
	logger    Logger
	processor OrderProcessor
}

func NewConsumer(cfg *Config, logger Logger, processor OrderProcessor) *Consumer {
	return &Consumer{
		config:    cfg,
		logger:    logger,
		processor: processor,
	}
}

func (c *Consumer) Init() error {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.config.Brokers,
		Topic:          c.config.Topic,
		CommitInterval: 0,
	})
	c.logger.Info("Kafka consumer initialized for topic: %s", c.config.Topic)
	return nil
}

func (c *Consumer) Run(ctx context.Context) {
	go c.Read(ctx)
}

func (c *Consumer) Stop() {
	if err := c.reader.Close(); err != nil {
		c.logger.Warn("Kafka consumer failed to close: %v", err)
	}
}

func (c *Consumer) Read(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Consumer context canceled, stopping...")
			return
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Error("Kafka read error: %v", err)
				continue
			}

			var order models.Order
			if err := json.Unmarshal(m.Value, &order); err != nil {
				c.logger.Error("Failed to unmarshal order: %v", err)
				continue
			}

			c.logger.Info("Received order ID: %d", order.ID)

			if err := c.processor.ProcessOrder(ctx, &order); err != nil {
				c.logger.Error("Failed to process order: %v", err)
			}
		}
	}
}
