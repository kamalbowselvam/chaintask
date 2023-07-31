package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)

type TaskService interface {
	GetTask(id int64) (domain.Task, error)
	CreateTask(context.Context, db.CreateTaskParams) (domain.Task,error)
	CreateUser(context.Context, db.CreateUserParams) (domain.User, error)
	GetUser(context.Context, string ) (domain.User, error)
	//SaveTask(domain.Task) (domain.Task, error)
	//DeleteTask(id int64) error
	//ListTask(page int64, offset int64) ([]domain.Task, error)
	//UpdateTask(task domain.Task) (domain.Task, error)
}