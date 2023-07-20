package db

import (
	"database/sql"
	"github.com/kamalbowselvam/chaintask/models"
)

type Store struct{
	db *sql.DB
	q *Queries
}

type Queries interface {
	CreateTask(Name string) (*models.Task, error)
	DeleteTask(Id int64) (error)
}