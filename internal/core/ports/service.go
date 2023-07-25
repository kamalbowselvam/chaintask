package ports

import "github.com/kamalbowselvam/chaintask/internal/core/domain"




type TaskService interface {
	GetTask(id int64) (domain.Task, error)
	CreateTask(name string, budget float64, user string) (domain.Task,error)
	//SaveTask(domain.Task) (domain.Task, error)
	//DeleteTask(id int64) error
	//ListTask(page int64, offset int64) ([]domain.Task, error)
	//UpdateTask(task domain.Task) (domain.Task, error)
}