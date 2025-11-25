package usecase

import (
	"food-story/pkg/exceptions"
	shareModel "food-story/shared/model"
	"log/slog"
	"sort"
)

func (i *Implement) PublishOrderToQueue(orderItems []*shareModel.OrderItems) error {

	if orderItems == nil {
		return exceptions.Error(exceptions.CodeBusiness, "order items is empty")
	}

	sort.Slice(orderItems, func(i, j int) bool {
		return orderItems[i].ID < orderItems[j].ID
	})

	for _, v := range orderItems {
		err := i.queue.PublishOrder(*v)
		if err != nil {
			slog.Error("failed to publish order to queue: ", " order items: ", err)
			continue
		}
	}

	return nil
}
