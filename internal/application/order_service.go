package application

import (
	"consumer-service/internal/models"
	"context"
	"github.com/pkg/errors"
)

type iRepository interface {
	SaveOrder(ctx context.Context, order *models.Order) error
}

type iCache interface {
	SaveOrder(ctx context.Context, order *models.Order) error
}

type OrderService struct {
	repo  iRepository
	cache iCache
}

func NewOrderService(repo iRepository, cache iCache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) Init() error {
	return nil
}

func (s *OrderService) Run(_ context.Context) {
}

func (s *OrderService) Stop() {
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
	if err := s.repo.SaveOrder(ctx, order); err != nil {
		return errors.Wrap(err, "repo.SaveOrder failed")
	}
	if err := s.cache.SaveOrder(ctx, order); err != nil {
		return errors.Wrap(err, "cache.SaveOrder failed")
	}
	return nil
}
