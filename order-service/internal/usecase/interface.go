package usecase

import (
	"context"
	"food-story/order-service/internal/adapter/cache"
	"food-story/order-service/internal/adapter/queue/producer"
	"food-story/pkg/exceptions"
	"food-story/shared/config"
	"github.com/google/uuid"

	"food-story/order-service/internal/adapter/repository"
	"food-story/order-service/internal/domain"
)

type Usecase interface {
	CreateOrder(ctx context.Context, sessionID uuid.UUID, items []domain.OrderItems) (result int64, customError *exceptions.CustomError)
	GetOrderByID(ctx context.Context, sessionID uuid.UUID) (result *domain.Order, customError *exceptions.CustomError)
	CreateOrderItems(ctx context.Context, sessionID uuid.UUID, items []domain.OrderItems) (customError *exceptions.CustomError)
	GetOrderItems(ctx context.Context, sessionID uuid.UUID) (result []*domain.OrderItems, customError *exceptions.CustomError)
	GetOderItemsByID(ctx context.Context, sessionID uuid.UUID, orderItemsID int64) (result *domain.OrderItems, customError *exceptions.CustomError)
	UpdateOrderItemsStatus(ctx context.Context, sessionID uuid.UUID, payload domain.OrderItemsStatus) (customError *exceptions.CustomError)
}

type Implement struct {
	config     config.Config
	repository repository.Implement
	cache      cache.RedisTableCacheInterface
	queue      producer.OrderProducer
}

func NewUsecase(config config.Config, repository repository.Implement, cache cache.RedisTableCacheInterface, queue producer.OrderProducer) *Implement {
	return &Implement{
		config,
		repository,
		cache,
		queue,
	}
}

var _ Usecase = (*Implement)(nil)
