package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"go.uber.org/zap"
)

type service struct {
	globalRepository   db.GlobalRepository
	policiesRepository authorization.PolicyManagementService
	logger             *zap.Logger
}

func NewTaskService(globalRepository db.GlobalRepository, policiesRepository authorization.PolicyManagementService, logger *zap.Logger) *service {
	return &service{
		globalRepository:   globalRepository,
		policiesRepository: policiesRepository,
		logger:             logger,
	}
}

func (srv *service) GetTask(ctx context.Context, id int64) (domain.Task, error) {
	task, err := srv.globalRepository.GetTask(context.Background(), id)
	return task, err
}

func (srv *service) CreateTask(ctx context.Context, arg db.CreateTaskParams) (domain.Task, error) {

	task, err := srv.globalRepository.CreateTask(context.Background(), arg)
	if err != nil {
		srv.logger.Fatal("Could not save the task in repository", zap.Error(err))

	}
	return task, err

}

func (srv *service) DeleteTask(ctx context.Context, id int64) error {
	err := srv.globalRepository.DeleteTask(context.Background(), id)
	if err != nil {
		srv.logger.Fatal("could not delete task in repository", zap.Error(err))
	}
	return err
}

func (srv *service) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	task, err := srv.globalRepository.UpdateTask(context.Background(), task)
	if err != nil {
		srv.logger.Fatal("could not update task in repository", zap.Error(err))
	}
	return task, err
}
