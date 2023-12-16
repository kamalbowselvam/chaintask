package service

import (
	"context"
	"fmt"

	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	"go.uber.org/zap"
)

type service struct {
	globalRepository   db.GlobalRepository
	policiesRepository authorization.PolicyManagementService
}

func NewTaskService(globalRepository db.GlobalRepository, policiesRepository authorization.PolicyManagementService) *service {
	return &service{
		globalRepository:   globalRepository,
		policiesRepository: policiesRepository,
	}
}

func (srv *service) GetTask(ctx context.Context, id int64) (domain.Task, error) {
	task, err := srv.globalRepository.GetTask(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), id)
	return task, err
}

func (srv *service) CreateTask(ctx context.Context, arg db.CreateTaskParams) (domain.Task, error) {
	logger_ := logger.FromCtx(ctx)
	task, err := srv.globalRepository.CreateTask(logger.WithCtx(context.Background(), logger_), arg)
	if err != nil {
		logger_.Error("Could not save the task in repository", zap.Error(err))
		return task, err
	}
	err = srv.policiesRepository.CreateTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)
	if err != nil {
		logger_.Error("could not create policy for task", zap.Int64("id", task.Id), zap.String(" due to ", err.Error()))
		err := srv.DeleteTask(ctx, task.Id)
		if err != nil {
			logger_.Error("could not delete task ", zap.Int64("id", task.Id), zap.String(" due to ", err.Error()))
		}
		return domain.Task{}, err
	}
	return task, nil

}

func (srv *service) DeleteTask(ctx context.Context, id int64) error {
	logger_ := logger.FromCtx(ctx)
	task, err := srv.globalRepository.GetTask(logger.WithCtx(context.Background(), logger_), id)
	if err != nil {
		return fmt.Errorf("trying to delete a task that does not exists %d", id)
	}
	err = srv.globalRepository.DeleteTask(logger.WithCtx(context.Background(), logger_), id)
	if err != nil {
		logger_.Error("could not delete task in repository", zap.Error(err))
	}
	err2 := srv.policiesRepository.RemoveTaskPolicies(id, task.ProjectId, task.CreatedBy, task.CompanyId)
	if err2 != nil {
		logger_.Error("could not delete policies linked to tasks", zap.Error(err2))

	}
	return err
}

func (srv *service) UpdateTask(ctx context.Context, task db.UpdateTaskParams) (domain.Task, error) {
	logger_ := logger.FromCtx(ctx)
	full_task, err := srv.globalRepository.UpdateTask(logger.WithCtx(context.Background(), logger_), task)
	if err != nil {
		logger_.Error("could not update task in repository", zap.Error(err))

	}
	return full_task, err
}
