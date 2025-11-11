package database

import (
	"context"
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

		return nil
	})

	return err
}
