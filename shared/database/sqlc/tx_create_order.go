package database

import (
	"context"
	"errors"
)

type TXCreateOrderParams struct {
	CreateOrder      CreateOrderParams
	CreateOrderItems []CreateOrderItemsParams
}

func (store *SQLStore) TXCreateOrder(ctx context.Context, arg TXCreateOrderParams) (int64, error) {

	var orderID int64
	err := store.execTx(ctx, func(q *Queries) error {
		orderIDRaw, orderError := q.CreateOrder(ctx, arg.CreateOrder)
		if orderError != nil {
			return orderError
		}
		orderID = orderIDRaw

		if len(arg.CreateOrderItems) > 0 {
			for index, _ := range arg.CreateOrderItems {
				arg.CreateOrderItems[index].OrderID = orderID
			}
			_, itemsError := q.CreateOrderItems(ctx, arg.CreateOrderItems)
			if itemsError != nil {
				return itemsError
			}
		}

		err := q.UpdateTablesStatusOrdered(ctx, arg.CreateOrder.TableID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	if orderID == 0 {
		return 0, errors.New("order id is 0")
	}

	return orderID, nil
}
