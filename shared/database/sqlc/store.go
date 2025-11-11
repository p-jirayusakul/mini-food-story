package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	TXCreateTableSession(ctx context.Context, arg CreateTableSessionParams) error
	TXCreateOrder(ctx context.Context, arg TXCreateOrderParams) (int64, error)
	TXSessionsExtension(ctx context.Context, arg TXSessionsExtensionParams) error
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
