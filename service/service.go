package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)

type TaskService interface {
	GetTask(context.Context, int64) (domain.Task, error)
	CreateTask(context.Context, db.CreateTaskParams) (domain.Task, error)
	DeleteTask(context.Context, int64) error
	UpdateTask(context.Context, db.UpdateTaskParams) (domain.Task, error)

	CreateUser(context.Context, db.CreateUserParams) (domain.User, error)
	GetUser(context.Context, string) (domain.User, error)
	CreateSession(context.Context, db.CreateSessionParams) (domain.Session, error)
	GetSession(context.Context, uuid.UUID) (domain.Session, error)
	DeleteUser(context.Context, string) error
	//FIXME Think about arguments that will go into this service
	//UpdateUser(context.Context, domain.User) (domain.User, error)

	CreateProject(context.Context, db.CreateProjectParam) (domain.Project, error)
	GetProject(context.Context, int64) (domain.Project, error)
	DeleteProject(context.Context, int64) error
	//FIXME Think about arguments that will go into this service
	//UpdateProject(context.Context, domain.Project) (domain.Project, error)

	CreateCompany(context.Context, db.CreateCompanyParams) (domain.Company, error)
}
