package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamalbowselvam/chaintask/domain"
)

type GlobalRepository interface {
	TaskRepository
	UserRepository
	ProjectRepository
	CompanyRepository
	ResourceRepository
	PayeeRepository
}

type UserRepository interface {
	CreateUser(context.Context, CreateUserParams) (domain.User, error)
	GetUser(context.Context, string) (domain.User, error)
	CreateSession(context.Context, CreateSessionParams) (domain.Session, error)
	GetSession(context.Context, uuid.UUID) (domain.Session, error)
	DeleteUser(context.Context, string) error
}

type TaskRepository interface {
	GetTask(context.Context, int64) (domain.Task, error)
	CreateTask(context.Context, CreateTaskParams) (domain.Task, error)
	GetTaskList(context.Context, []int64) ([]domain.Task, error)
	GetTaskListByProject(context.Context, int64) ([]domain.Task, error)
	DeleteTask(context.Context, int64) error
	UpdateTask(context.Context, UpdateTaskParams) (domain.Task, error)
	DeleteTasksLinkedToProject(context.Context, int64) error
}

type ProjectRepository interface {
	CreateProject(context.Context, CreateProjectParam) (domain.Project, error)
	GetClientAndResponsibleByProject(context.Context, int64) (string, string, error)
	DeleteProject(context.Context, int64) error
	GetProject(context.Context, int64) (domain.Project, error)
}

type CompanyRepository interface {
	CreateCompany(context.Context, CreateCompanyParams) (domain.Company, error)
}

type ResourceRepository interface {
	CreateResource(context.Context, CreateResourceParams) (domain.Resource, error)
}


type PayeeRepository interface {
	CreatePayee(context.Context, CreatePayeeParams) (domain.Payee, error)
}

var _ GlobalRepository = (*Queries)(nil)
