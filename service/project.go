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
		srv.policiesRepository.CreateProjectPolicies(project.Id, project.Client, project.Responsible)
		// FIXME If policies fails, then the project should be deleted
	}
	return project, err
}
