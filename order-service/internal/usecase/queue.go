package usecase

import (
	"errors"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"log/slog"
	"sort"
)

func (i *Implement) PublishOrderToQueue(orderItems []*domain.OrderItems) *exceptions.CustomError {

	if orderItems == nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: errors.New("order items is empty"),
		}
	}

	sort.Slice(orderItems, func(i, j int) bool {
		return orderItems[i].ID < orderItems[j].ID
	})

	for _, v := range orderItems {
		err := i.queue.PublishOrder(*v)
		if err != nil {
			slog.Error("failed to publish order to queue: ", err, " order items: ", v)
			continue
		}
	}

	return nil
}
