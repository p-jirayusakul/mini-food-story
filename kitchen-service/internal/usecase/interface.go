package usecase

import (
	"context"
	"food-story/kitchen-service/internal/adapter/repository"
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/shared/config"
)

type Usecase interface {
	UpdateOrderItemsStatus(ctx context.Context, payload domain.OrderItemsStatus) (customError *exceptions.CustomError)
	SearchOrderItems(ctx context.Context, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, customError *exceptions.CustomError)
	GetOrderItems(ctx context.Context, orderID int64) (result []*domain.OrderItems, customError *exceptions.CustomError)
	GetOrderItemsByID(ctx context.Context, orderID, orderItemsID int64) (result *domain.OrderItems, customError *exceptions.CustomError)
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
