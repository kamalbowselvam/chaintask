package db

import (
	"database/sql"

	"go.uber.org/zap"
)

func New(db *sql.DB) *Queries {
	return &Queries{
		db:     db,
	}
}

type Queries struct {
	db     *sql.DB
	logger *zap.Logger
}
