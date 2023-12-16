package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	"go.uber.org/zap"
)

func (srv *service) CreateProject(ctx context.Context, arg db.CreateProjectParam) (domain.Project, error) {
	logger_ := logger.FromCtx(ctx)
	project, err := srv.globalRepository.CreateProject(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), arg)
	if err != nil {

		logger_.Error("could not create project due to", zap.Error(err))

		return domain.Project{}, err
	}
	project.CompletionPercentage = 0
	project.Budget = 0
	if err == nil {
		err = srv.policiesRepository.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
		if err != nil {
			err = srv.DeleteProject(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), project.Id)
			if err != nil {
				logger_.Error("Could not delete project, invalid policies may be created", zap.Error(err))
			}
			return domain.Project{}, err
		}
	}
	return project, err
}

func (srv *service) DeleteProject(ctx context.Context, id int64) error {
	logger_ := logger.FromCtx(ctx)
	project, err := srv.GetProject(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), id)
	if err != nil {
		return nil
	}
	srv.policiesRepository.RemoveProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
	for _, task := range project.Tasks {
		srv.policiesRepository.RemoveTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, project.CompanyId)
	}
	err = srv.globalRepository.DeleteTasksLinkedToProject(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), id)
	if err != nil {
		logger_.Error("could not tasks linked to project ", zap.Int64(" id ", id), zap.String("due to ", err.Error()))
	}
	err = srv.globalRepository.DeleteProject(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), id)
	if err != nil {
		logger_.Error("could not deleted project ", zap.Int64(" id ", id), zap.String("due to ", err.Error()))
		return err
	}
	return nil
}

func (srv *service) GetProject(ctx context.Context, id int64) (domain.Project, error) {
	project, err := srv.globalRepository.GetProject(ctx, id)
	if err != nil {
		return domain.Project{}, nil
	}
	tasks, err := srv.globalRepository.GetTaskListByProject(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), id)
	if err != nil {
		return domain.Project{}, nil
	}
	project.Tasks = tasks
	return project, nil
}
