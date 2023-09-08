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
	UpdateTask(context.Context, domain.Task) (domain.Task,error)

	CreateUser(context.Context, db.CreateUserParams) (domain.User, error)
	GetUser(context.Context, string) (domain.User, error)
	CreateSession(context.Context, db.CreateSessionParams) (domain.Session,error)
	GetSession(context.Context, uuid.UUID) (domain.Session, error)
	
	
	// DeleteUser
	// UpdateUser

	// Create Project
	CreateProject(context.Context, db.CreateProjectParam) (domain.Project, error)
	// Get Project
	// Delete Project
	// Update Project

}
