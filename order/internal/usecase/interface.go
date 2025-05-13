package usecase

import (
	"context"
	"food-story/order/internal/adapter/cache"
	"food-story/pkg/exceptions"
	"food-story/shared/config"
	"github.com/google/uuid"

	//"food-story/order/internal/adapter/cache"
	"food-story/order/internal/adapter/repository"
	"food-story/order/internal/domain"
)

type Usecase interface {
	CreateOrder(ctx context.Context, sessionID uuid.UUID) (result int64, customError *exceptions.CustomError)
	GetOrderByID(ctx context.Context, id int64) (result *domain.Order, customError *exceptions.CustomError)
	UpdateOrderStatus(ctx context.Context, sessionID uuid.UUID, payload domain.OrderStatus) (customError *exceptions.CustomError)
	CreateOrderItems(ctx context.Context, sessionID uuid.UUID, items []domain.OrderItems) (customError *exceptions.CustomError)
	GetOrderItems(ctx context.Context, orderID int64) (result []*domain.OrderItems, customError *exceptions.CustomError)
	GetOderItemsByID(ctx context.Context, orderID, orderItemsID int64) (result *domain.OrderItems, customError *exceptions.CustomError)
	UpdateOrderItemsStatus(ctx context.Context, sessionID uuid.UUID, payload domain.OrderItemsStatus) (customError *exceptions.CustomError)
}

type Implement struct {
	config     config.Config
	repository repository.Implement
	cache      cache.RedisTableCacheInterface
}

func NewUsecase(config config.Config, repository repository.Implement, cache cache.RedisTableCacheInterface) *Implement {
	return &Implement{
		config,
		repository,
		cache,
	}
}

var _ Usecase = (*Implement)(nil)
