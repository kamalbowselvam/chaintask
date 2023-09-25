package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"go.uber.org/zap"
)

func (srv *service) CreateProject(ctx context.Context, arg db.CreateProjectParam) (domain.Project, error) {
	project, err := srv.globalRepository.CreateProject(context.Background(), arg)
	if err != nil {

		srv.logger.Fatal("could not create project due to", zap.Error(err))

		return domain.Project{}, err
	}
	project.CompletionPercentage = 0
	project.Budget = 0
	if err == nil {
		err = srv.policiesRepository.CreateProjectPolicies(project.Id, project.Client, project.Responsible)
		if err != nil {
			err = srv.DeleteProject(ctx, project.Id)
			if err != nil {
				srv.logger.Fatal("Could not delete project, invalid policies may be created", zap.Error(err))
			}
			return domain.Project{}, err
		}
	}
	return project, err
}

func (srv *service) DeleteProject(ctx context.Context, id int64) error {
	project, err := srv.GetProject(ctx, id)
	if err != nil {
		return nil
	}
	srv.policiesRepository.RemoveProjectPolicies(project.Id, project.Client, project.Responsible)
	for _, task := range project.Tasks{
		srv.policiesRepository.RemoveTaskPolicies(task.Id, task.ProjectId, task.CreatedBy)
	}
	err = srv.globalRepository.DeleteTasksLinkedToProject(ctx, id)
	if err != nil {
		srv.logger.Fatal("could not tasks linked to project ", zap.Int64(" id ", id), zap.String("due to ", err.Error()))
	}
	err = srv.globalRepository.DeleteProject(ctx, id)
	if err != nil {
		srv.logger.Fatal("could not deleted project ", zap.Int64(" id ", id), zap.String("due to ", err.Error()))
		return err
	}
	return nil
}

func (srv *service) GetProject(ctx context.Context, id int64) (domain.Project, error) {
	project, err := srv.globalRepository.GetProject(ctx, id)
	if err != nil {
		return domain.Project{}, nil
	}
	tasks, err := srv.globalRepository.GetTaskListByProject(ctx, id)
	if err != nil {
		return domain.Project{}, nil
	}
	project.Tasks = tasks
	return project, nil
}
