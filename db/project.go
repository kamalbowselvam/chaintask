package db

import (
	"context"

	"github.com/kamalbowselvam/chaintask/domain"
)

const createProject = `INSERT INTO projects (
	projectname,
	created_by,
	location,
	address,
	responsible,
	client
  ) VALUES (
	$1, $2, $3, $4, $5, $6
  )
  RETURNING id, projectname, created_on, created_by, location, address, responsible, client;`

type CreateProjectParam struct {
	ProjectName string          `json:"projectname"`
	CreatedBy   string          `json:"createdBy"`
	Location    domain.Location `json:"location"`
	Address     string          `json:"address"`
	Responsible string          `json:"responsible"`
	Client      string          `json:"client"`
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParam) (domain.Project, error) {
	q.logger.Info("saving projects")

	row := q.db.QueryRowContext(ctx, createProject, arg.ProjectName, arg.CreatedBy, Point{arg.Location[0], arg.Location[1]}, arg.Address, arg.Responsible, arg.Client)
	var i domain.Project
	var p Point
	err := row.Scan(
		&i.Id,
		&i.Projectname,
		&i.CreatedOn,
		&i.CreatedBy,
		&p,
		&i.Address,
		&i.Responsible,
		&i.Client,
	)
	i.Location = domain.Location{p[0], p[1]}
	return i, err

}

const getClientAndResponsibleByProject = `SELECT client, responsible from projects where id=$1`

func (q *Queries) GetClientAndResponsibleByProject(ctx context.Context, projectId int64) (string, string, error) {
	row := q.db.QueryRowContext(ctx, getClientAndResponsibleByProject, projectId)
	client := ""
	responsible := ""
	err := row.Scan(
		client,
		responsible,
	)
	return client, responsible, err
}

const deleteProject =  `delete from projects where id=$1`
func (q *Queries) DeleteProject(ctx context.Context, projectId int64) (error){
	_, err := q.db.ExecContext(ctx, deleteProject, projectId)
	return err
}

const getProject = `select projectname, created_on, created_by, location, address, responsible, client FROM projects where id = $1`
func (q *Queries) GetProject(ctx context.Context, projectId int64) (domain.Project, error){
	row := q.db.QueryRowContext(ctx, getProject, projectId)
	var i domain.Project
	var p Point
	err := row.Scan(
		&i.Id,
		&i.Projectname,
		&i.CreatedOn,
		&i.CreatedBy,
		&p,
		&i.Address,
		&i.Responsible,
		&i.Client,
	)
	i.Location = domain.Location{p[0], p[1]}
	return i, err
}