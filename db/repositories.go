package db

import (
	"context"

	"github.com/kamalbowselvam/chaintask/domain"
)



type GlobalRepository interface {
	TaskRepository
	UserRepository
	ProjectRepository
}

type UserRepository interface {
	CreateUser(context.Context, CreateUserParams) (domain.User, error)
	GetUser(context.Context, string)(domain.User, error)
}

type TaskRepository interface {
	
	GetTask(context.Context, int64) (domain.Task, error)
	CreateTask(context.Context, CreateTaskParams) (domain.Task, error)
	GetTaskList(context.Context, []int64) ([]domain.Task, error)
	GetTaskListByProject(context.Context, int64) ([]domain.Task, error)
	DeleteTask(context.Context, int64) error
	UpdateTask(context.Context, domain.Task) (domain.Task, error)
}

type ProjectRepository interface {
	CreateProject(context.Context, CreateProjectParam) (domain.Project, error)
}





