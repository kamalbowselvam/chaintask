package db

import (
	"database/sql"

	"go.uber.org/zap"
)

func New(db *sql.DB, logger *zap.Logger) *Queries {
	return &Queries{
		db:     db,
		logger: logger,

	}
}

type Queries struct {
	db     *sql.DB
	logger *zap.Logger
}
