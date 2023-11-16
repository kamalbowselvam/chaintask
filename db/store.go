package db

import (
	"database/sql"

	"github.com/kamalbowselvam/chaintask/logger"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	GlobalRepository
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *sql.DB) Store {

	logger.Info("Starting the Database connection")
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
