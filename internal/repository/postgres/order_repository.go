package postgres

import (
	"consumer-service/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type PostgresRepository struct {
	conn   *pgx.Conn
	config *Config
	logger Logger
}

func NewPostgresRepository(cfg *Config, logger Logger) *PostgresRepository {
	return &PostgresRepository{
		config: cfg,
		logger: logger,
	}
}

func (r *PostgresRepository) Init() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		r.config.Host,
		r.config.Port,
		r.config.Username,
		r.config.DBName,
		r.config.Password,
		r.config.SSLMode,
	)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}

	if err = conn.Ping(context.Background()); err != nil {
		return err
	}

	r.conn = conn
	r.logger.Info("Connected to Postgres")
	return nil
}

func (r *PostgresRepository) Run(_ context.Context) {}

func (r *PostgresRepository) Stop() {
	if err := r.conn.Close(context.Background()); err != nil {
		r.logger.Warn("Failed to close Postgres connection: %v", err)
	}
}

func (r *PostgresRepository) SaveOrder(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (id, user_id, product_name, quantity, price, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.conn.Exec(ctx, query,
		order.ID,
		order.UserID,
		order.ProductName,
		order.Quantity,
		order.Price,
		order.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	r.logger.Info("Order saved to Postgres: ID %d", order.ID)
	return nil
}
