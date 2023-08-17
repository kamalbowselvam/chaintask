package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)

type TaskService interface {
	GetTask(context.Context, db.GetTaskParams) (domain.Task, error)
	CreateTask(context.Context, db.CreateTaskParams) (domain.Task,error)
	// DeleteTask(context.Context, int64) error TODO
	// UpdateTask(context.Contect, db.UpdateTaskParams) (domain.Task,error)
	
	
	CreateUser(context.Context, db.CreateUserParams) (domain.User, error)
	GetUser(context.Context, string ) (domain.User, error)
	// DeleteUser
	// UpdateUser 


	// Create Project 
	// Get Project 
	// Delete Project 
	// Update Project 

}