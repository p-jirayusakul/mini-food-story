package usecase

import (
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"sort"
)

func (i *Implement) PublishOrderToQueue(orderItems []*domain.OrderItems) *exceptions.CustomError {
	// Sort by ID ascending
	sort.Slice(orderItems, func(i, j int) bool {
		return orderItems[i].ID < orderItems[j].ID
	})

	// public message to kafka
	for _, v := range orderItems {
		err := i.queue.PublishOrder(*v)
		if err != nil {
			return &exceptions.CustomError{
				Status: exceptions.ERRREPOSITORY,
				Errors: err,
			}
		}
	}

	return nil
}
