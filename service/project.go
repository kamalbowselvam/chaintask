package service

import (
	"context"
	"log"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)


func (srv *service) CreateProject(ctx context.Context, arg db.CreateProjectParam) (domain.Project, error) {
	project, err := srv.globalRepository.CreateProject(context.Background(), arg)
	if err != nil {
		log.Fatalf("could not create project due to %s", err)
		return domain.Project{}, err
	}
	project.CompletionPercentage = 0
	project.Budget = 0
	if err != nil {
		srv.policiesRepository.CreateProjectPolicies(project.Id, project.Client, project.Responsible)
	}
	return project, err
}
