package db
import "database/sql"


func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db *sql.DB
}
