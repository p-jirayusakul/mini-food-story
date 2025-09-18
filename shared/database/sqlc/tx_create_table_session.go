package database

import (
	"context"
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
