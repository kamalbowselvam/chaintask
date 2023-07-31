package service

import (
	"context"
	"log"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)


type service struct {
	taskRepository db.TaskRepository
}


func NewTaskService(taskRepository db.TaskRepository) *service {
	return &service {
		taskRepository: taskRepository,
	}
}


func (srv *service) GetTask(id int64) (domain.Task, error){
	task ,err := srv.taskRepository.GetTask(context.Background(),id)
	return task, err
}


func (srv *service) CreateTask(ctx context.Context, arg db.CreateTaskParams) (domain.Task,error){

	//task := domain.NewTask(name, budget, user)
	task, err := srv.taskRepository.CreateTask(context.Background(),arg)
	if err != nil {
		log.Fatal("Could not save the task in repository", err.Error())

	}

	return task, err

}