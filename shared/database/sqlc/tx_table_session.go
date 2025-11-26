package database

import (
	"context"
	"errors"
	"food-story/pkg/utils"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (store *SQLStore) TXCreateTableSession(ctx context.Context, arg CreateTableSessionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {
		err := q.CreateTableSession(ctx, arg)
		if err != nil {
			return err
		}

		err = q.UpdateTablesStatusWaitToOrder(ctx, arg.TableID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

type TXSessionsExtensionParams struct {
	TableID                int64
	NewOrderItemsID        int64
	ProductID              int64
	RequestedMinutes       int64
	ReasonCode             string
	ExpiresAt              time.Time
	CreateSessionExtension CreateSessionExtensionParams
}

func (store *SQLStore) TXSessionsExtension(ctx context.Context, arg TXSessionsExtensionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		sessionID, err := q.GetSessionIDByTableID(ctx, arg.TableID)
		if err != nil {
			return err
		}
		arg.CreateSessionExtension.SessionID = sessionID

		_, err = q.CreateSessionExtension(ctx, arg.CreateSessionExtension)
		if err != nil {
			return err
		}

		expiresAt := pgtype.Timestamptz{
			Time:  arg.ExpiresAt,
			Valid: true,
		}
		err = q.UpdateSessionExpireBySessionID(ctx, UpdateSessionExpireBySessionIDParams{
			RequestedMinutes: int32(arg.RequestedMinutes),
			LastReasonCode:   arg.ReasonCode,
			ExpiresAt:        expiresAt,
			Sessionid:        sessionID,
		})
		if err != nil {
			return err
		}

		product, err := q.GetProductByID(ctx, arg.ProductID)
		if err != nil {
			return err
		}

		if product == nil {
			return errors.New("product is null")
		}

		statusServedID, err := q.GetOrderStatusServed(ctx)
		if err != nil {
			return err
		}

		currentTime, err := q.GetTimeNow(ctx)
		if err != nil {
			return err
		}

		isFree, err := q.IsSessionExtensionModeFree(ctx, arg.CreateSessionExtension.ModeID.Int64)
		if err != nil {
			return err
		}

		var price pgtype.Numeric
		if isFree {
			price = utils.Float64ToPgNumeric(0)
		} else {
			productPrice, fErr := product.Price.Float64Value()
			if fErr != nil {
				return fErr
			}
			price = utils.Float64ToPgNumeric(productPrice.Float64)
		}

		orderID, err := q.GetOrderIDBySessionID(ctx, sessionID)
		if err != nil {
			return err
		}

		createOrderItemsParams := CreateOrderItemsPerRowParams{
			ID:              arg.NewOrderItemsID,
			OrderID:         orderID,
			ProductID:       product.ID,
			StatusID:        statusServedID,
			ProductName:     product.Name,
			ProductNameEn:   product.NameEn,
			Price:           price,
			Quantity:        1,
			CreatedAt:       currentTime,
			ProductImageUrl: product.ImageUrl,
			IsVisible:       false,
		}
		err = q.CreateOrderItemsPerRow(ctx, createOrderItemsParams)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
