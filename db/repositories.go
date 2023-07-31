package db

import (
	"context"

	"github.com/kamalbowselvam/chaintask/domain"
)



type GlobalRepository interface {
	TaskRepository
	UserRepository
}

type UserRepository interface {
	CreateUser(context.Context, domain.User) (domain.User, error)
}

type TaskRepository interface {
	
	GetTask(context.Context, int64) (domain.Task, error)
	CreateTask(context.Context, CreateTaskParams) (domain.Task, error)
	GetTaskList(context.Context, []int64) ([]domain.Task, error)
	DeleteTask(context.Context, int64) error
	UpdateTask(context.Context, domain.Task) (domain.Task, error)
}





