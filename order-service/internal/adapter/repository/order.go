package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (i *Implement) CreateOrder(ctx context.Context, order domain.CreateOrder) (orderID int64, err error) {

	orderItems, err := i.buildPayloadOrderItems(ctx, order.OrderItems)
	if err != nil {
		return 0, err
	}

	orderNumber, err := i.getOrderSequence(ctx)
	if err != nil {
		return 0, err
	}

	orderID, err = i.repository.TXCreateOrder(ctx, database.TXCreateOrderParams{
		CreateOrderItems: orderItems,
		CreateOrder: database.CreateOrderParams{
			ID:          i.snowflakeID.Generate(),
			OrderNumber: orderNumber,
			SessionID:   utils.UUIDToPgUUID(order.SessionID),
			TableID:     order.TableID,
		},
	})
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to create order", err)
	}

	return orderID, nil
}
func (i *Implement) GetOrderByID(ctx context.Context, id int64) (result *domain.Order, err error) {
	order, err := i.repository.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderNotFound.Error())
		}
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to check order exists", err)
	}

	if order == nil {
		return nil, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderNotFound.Error())
	}

	return &domain.Order{
		ID:           order.ID,
		TableID:      order.TableID,
		TableNumber:  order.TableNumber,
		StatusID:     order.StatusID,
		StatusName:   order.StatusName,
		StatusNameEN: order.StatusNameEN,
		StatusCode:   order.StatusCode,
	}, nil
}

func (i *Implement) IsOrderExist(ctx context.Context, id int64) (err error) {
	isExist, err := i.repository.IsOrderExist(ctx, id)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to check order exists", err)
	}

	if !isExist {
		return exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderNotFound.Error())
	}

	return nil
}

func (i *Implement) IsOrderWithItemsExists(ctx context.Context, orderID, orderItemsID int64) (err error) {
	isExist, err := i.repository.IsOrderWithItemsExists(ctx, database.IsOrderWithItemsExistsParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to check order item exists", err)
	}

	if !isExist {
		return exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderItemsNotFound.Error())
	}

	return nil
}

func (i *Implement) IsOrderItemsNotFinal(ctx context.Context, orderID int64) (bool, error) {
	isExist, err := i.repository.IsOrderItemsNotFinal(ctx, orderID)
	if err != nil {
		return false, exceptions.Errorf(exceptions.CodeRepository, "failed to check order items not final", err)
	}
	return isExist, nil
}

func (i *Implement) GetTableIDByOrderID(ctx context.Context, orderID int64) (int64, error) {
	tableID, err := i.repository.GetTableIDByOrderID(ctx, orderID)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to get table id by order id", err)
	}
	return tableID, nil
}

func (i *Implement) getOrderSequence(ctx context.Context) (string, error) {

	currentTimeDB, err := i.repository.GetTimeNow(ctx)
	if err != nil {
		return "", exceptions.Errorf(exceptions.CodeRepository, "failed to get current time", err)
	}

	if !currentTimeDB.Valid {
		return "", exceptions.Errorf(exceptions.CodeRepository, "validation failed: current time is not valid", err)
	}

	currentLocation, err := time.LoadLocation(i.config.TimeZone)
	if err != nil {
		return "", exceptions.Errorf(exceptions.CodeSystem, "failed to load time zone", err)
	}

	currentTime := currentTimeDB.Time.In(currentLocation)
	sequence, err := i.repository.GetOrderSequence(ctx, pgtype.Date{
		Time:  currentTime,
		Valid: true,
	})
	if err != nil {
		return "", exceptions.Errorf(exceptions.CodeRepository, "failed to get or create order sequence", err)
	}

	return fmt.Sprintf("FS-%s-%04d", currentTime.In(currentLocation).Format("20060102"), sequence), nil
}

func (i *Implement) GetSessionIDByOrderID(ctx context.Context, orderID int64) (result uuid.UUID, err error) {
	sessionIDData, err := i.repository.GetSessionIDByOrderID(ctx, orderID)
	if err != nil {
		return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to get table by order id", err)
	}

	sessionIDString := sessionIDData.String()
	sessionID, err := uuid.Parse(sessionIDString)
	if err != nil {
		return uuid.Nil, exceptions.Errorf(exceptions.CodeSystem, "failed to get session by order id", err)
	}

	return sessionID, nil
}
