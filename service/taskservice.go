package service

import (
	"context"
	"log"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)

type service struct {
	taskRepository db.GlobalRepository
}

func NewTaskService(taskRepository db.GlobalRepository) *service {
	return &service{
		taskRepository: taskRepository,
	}
}

func (srv *service) GetTask(ctx context.Context, id int64) (domain.Task, error) {
	task, err := srv.taskRepository.GetTask(context.Background(), id)
	return task, err
}

func (srv *service) CreateTask(ctx context.Context, arg db.CreateTaskParams) (domain.Task, error) {

	//task := domain.NewTask(name, budget, user)
	task, err := srv.taskRepository.CreateTask(context.Background(), arg)
	if err != nil {
		log.Fatal("Could not save the task in repository", err.Error())

	}

	return task, err

}

func (srv *service) DeleteTask(ctx context.Context, id int64) error {
	err := srv.taskRepository.DeleteTask(context.Background(), id)
	if err != nil {
		log.Fatalf("could not delete task in repository %s", err.Error())
	}
	return err
}

func (srv *service) CreateUser(ctx context.Context, arg db.CreateUserParams) (domain.User, error) {

	user, err := srv.taskRepository.CreateUser(context.Background(), arg)

	return user, err

}

func (srv *service) GetUser(ctx context.Context, username string) (domain.User, error) {
	user, err := srv.taskRepository.GetUser(context.Background(), username)

	return user, err
}
