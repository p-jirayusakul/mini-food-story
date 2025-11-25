package usecase

import (
	"context"
	"food-story/order-service/internal/adapter/cache"
	"food-story/order-service/internal/adapter/queue/producer"
	"food-story/shared/config"
	shareModel "food-story/shared/model"

	"github.com/google/uuid"

	"food-story/order-service/internal/adapter/repository"
	"food-story/order-service/internal/domain"
)

type Usecase interface {
	CreateOrder(ctx context.Context, sessionID uuid.UUID, items []shareModel.OrderItems) (result int64, err error)
	GetOrderByID(ctx context.Context, sessionID uuid.UUID) (result *domain.Order, err error)
	CreateOrderItems(ctx context.Context, sessionID uuid.UUID, items []shareModel.OrderItems) (err error)
	GetCurrentOrderItems(ctx context.Context, sessionID uuid.UUID, pageNumber, pageSize int64) (result domain.SearchCurrentOrderItemsResult, err error)
	GetCurrentOrderItemsByID(ctx context.Context, sessionID uuid.UUID, orderItemsID int64) (result *domain.CurrentOrderItems, err error)
	UpdateOrderItemsStatus(ctx context.Context, sessionID uuid.UUID, payload shareModel.OrderItemsStatus) (err error)
	GetOrderItems(ctx context.Context, orderID int64, pageNumber, pageSize int64) (result domain.SearchCurrentOrderItemsResult, err error)
	UpdateOrderItemsStatusByID(ctx context.Context, payload shareModel.OrderItemsStatus) (err error)
	SearchOrderItemsIncomplete(ctx context.Context, orderID int64, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, err error)
	IsSessionValid(sessionID uuid.UUID) error
	GetSessionIDByOrderID(ctx context.Context, orderID int64) (result uuid.UUID, err error)
	GetSessionIDByTableID(ctx context.Context, tableID int64) (result uuid.UUID, err error)
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
