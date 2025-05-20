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
	GetCurrentOrderItems(ctx context.Context, sessionID uuid.UUID) (result []*domain.CurrentOrderItems, customError *exceptions.CustomError)
	GetCurrentOderItemsByID(ctx context.Context, sessionID uuid.UUID, orderItemsID int64) (result *domain.CurrentOrderItems, customError *exceptions.CustomError)
	UpdateOrderItemsStatus(ctx context.Context, sessionID uuid.UUID, payload domain.OrderItemsStatus) (customError *exceptions.CustomError)
	SearchOrderItemsIncomplete(ctx context.Context, orderID int64, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, customError *exceptions.CustomError)
}

type Implement struct {
	config     config.Config
	repository repository.Implement
	cache      cache.RedisTableCacheInterface
	queue      producer.QueueProducerInterface
}

func NewUsecase(config config.Config, repository repository.Implement, cache cache.RedisTableCacheInterface, queue producer.QueueProducerInterface) *Implement {
	return &Implement{
		config,
		repository,
		cache,
		queue,
	}
}

var _ Usecase = (*Implement)(nil)
