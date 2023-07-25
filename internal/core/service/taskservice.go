package service

import (
	"context"
	"log"

	"github.com/kamalbowselvam/chaintask/internal/core/domain"
	"github.com/kamalbowselvam/chaintask/internal/core/ports"
)


type service struct {
	taskRepository ports.TaskRepository
}


func NewTaskService(taskRepository ports.TaskRepository) *service {
	return &service {
		taskRepository: taskRepository,
	}
}


func (srv *service) GetTask(id int64) (domain.Task, error){
	task ,err := srv.taskRepository.GetTask(context.Background(),id)
	return task, err
}


func (srv *service) CreateTask(name string, budget float64, user string) (domain.Task,error){

	task := domain.NewTask(name, budget, user)
	task, err := srv.taskRepository.SaveTask(context.Background(),task)
	if err != nil {
		log.Fatal("Could not save the task in repository", err.Error())

	}

	return task, err

}