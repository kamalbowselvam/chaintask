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

func (srv *service) CreateUser(ctx context.Context, arg db.CreateUserParams) (domain.User, error) {

	user, err := srv.globalRepository.CreateUser(context.Background(), arg)

	return user, err

}

func (srv *service) GetUser(ctx context.Context, username string) (domain.User, error) {
	user, err := srv.globalRepository.GetUser(context.Background(), username)

	return user, err
}

func (srv *service) CreateProject(ctx context.Context, arg db.CreateProjectParam) (domain.Project, error) {
	project, err := srv.globalRepository.CreateProject(context.Background(), arg)
	if err != nil {
		log.Fatalf("could not create project due to %s", err)
		return domain.Project{}, err
	}


	// shouldn't the task be empty when we create the project for first time ? 
	//tasks, err := srv.globalRepository.GetTaskListByProject(context.Background(), project.Id)
	//if err != nil {
	//	return project, err
	//}
	//project.Tasks = tasks
	// FIXME Compute global budget and completion stage
	//n := len(tasks)
	//if n > 0 {
	project.CompletionPercentage = 0
	//}
	project.Budget = 0
	return project, err
}
