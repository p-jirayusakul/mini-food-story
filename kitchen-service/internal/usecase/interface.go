package usecase

import (
	"context"
	"food-story/kitchen-service/internal/adapter/repository"
	"food-story/kitchen-service/internal/domain"
	"food-story/shared/config"
	shareModel "food-story/shared/model"
)

type Usecase interface {
	UpdateOrderItemsStatus(ctx context.Context, payload shareModel.OrderItemsStatus) (err error)
	UpdateOrderItemsStatusServed(ctx context.Context, payload shareModel.OrderItemsStatus) (err error)
	SearchOrderItems(ctx context.Context, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, err error)
	GetOrderItems(ctx context.Context, orderID int64, search domain.SearchOrderItems) (result domain.SearchOrderItemsResult, err error)
	GetOrderItemsByID(ctx context.Context, orderID, orderItemsID int64) (result *shareModel.OrderItems, err error)
}

type Implement struct {
	config     config.Config
	repository repository.Implement
}

func NewUsecase(config config.Config, repository repository.Implement) *Implement {
	return &Implement{
		config,
		repository,
	}
}

var _ Usecase = (*Implement)(nil)
