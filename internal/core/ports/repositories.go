package ports

import (
	"context"

	"github.com/kamalbowselvam/chaintask/internal/core/domain"
)

type TaskRepository interface {
	GetTask(context.Context, int64) (domain.Task, error)
	SaveTask(context.Context, domain.Task) (domain.Task, error)
	GetTaskList(context.Context, []int64) ([]domain.Task, error)
	DeleteTask(context.Context, int64) error
	UpdateTask(context.Context, domain.Task) (domain.Task, error)
}

type UserRepository interface {
	CreateUser(context.Context, domain.User) (domain.UserDetail, error)
	GetUser(context.Context, string) (domain.User, error)
}

type Storage interface {
	UserRepository
	TaskRepository
}
