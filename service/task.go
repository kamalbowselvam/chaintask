package service

import (
	"context"
	"log"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)

type service struct {
	globalRepository db.GlobalRepository
}

func NewTaskService(globalRepository db.GlobalRepository) *service {
	return &service{
		globalRepository: globalRepository,
	}
}

func (srv *service) GetTask(ctx context.Context, id int64) (domain.Task, error) {
	task, err := srv.globalRepository.GetTask(context.Background(), id)
	return task, err
}

func (srv *service) CreateTask(ctx context.Context, arg db.CreateTaskParams) (domain.Task, error) {

	task, err := srv.globalRepository.CreateTask(context.Background(), arg)
	if err != nil {
		log.Fatal("Could not save the task in repository", err.Error())

	}

	return task, err

}

func (srv *service) DeleteTask(ctx context.Context, id int64) error {
	err := srv.globalRepository.DeleteTask(context.Background(), id)
	if err != nil {
		log.Fatalf("could not delete task in repository %s", err.Error())
	}
	return err
}

func (srv *service) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	task, err := srv.globalRepository.UpdateTask(context.Background(), task)
	if err != nil {
		log.Fatalf("could not update task in repository %s", err.Error())
	}
	return task, err
}

