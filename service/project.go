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
		log.Printf("could not create project due to %s", err)
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