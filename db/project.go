package db

import (
	"context"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	"go.uber.org/zap"
)

const createProject = `INSERT INTO projects (
	projectname,
	created_by,
	longitude,
	latitude,
	address,
	responsible,
	client,
	company_id
  ) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8
  )
  RETURNING id, projectname, created_at, created_by, longitude, latitude, address, responsible, client, company_id;`

type CreateProjectParam struct {
	ProjectName string  `json:"project_name" binding:"required"`
	CreatedBy   string  `swaggerignore:"true"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Address     string  `json:"address"`
	Responsible string  `json:"responsible" binding:"required"`
	Client      string  `json:"client" binding:"required"`
	CompanyId   int64
}

type CompanyParam struct {
	CompanyId int64 `uri:"companyId" binding:"required,min=1"`
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParam) (domain.Project, error) {
	logger_ := logger.FromCtx(ctx)
	logger_.Info("saving projects", zap.Any("projects", arg))

	row := q.db.QueryRowContext(ctx, createProject, arg.ProjectName, arg.CreatedBy, arg.Longitude, arg.Latitude, arg.Address, arg.Responsible, arg.Client, arg.CompanyId)
	var i domain.Project
	err := row.Scan(
		&i.Id,
		&i.Projectname,
		&i.CreatedOn,
		&i.CreatedBy,
		&i.Longitude,
		&i.Latitude,
		&i.Address,
		&i.Responsible,
		&i.Client,
		&i.CompanyId,
	)
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

const deleteProject = `delete from projects where id=$1`

func (q *Queries) DeleteProject(ctx context.Context, projectId int64) error {
	_, err := q.db.ExecContext(ctx, deleteProject, projectId)
	return err
}

const getProject = `select projectname, created_at, created_by, longitude, latitude, address, responsible, client, company_id FROM projects where id = $1`

func (q *Queries) GetProject(ctx context.Context, projectId int64) (domain.Project, error) {
	row := q.db.QueryRowContext(ctx, getProject, projectId)
	var i domain.Project
	err := row.Scan(
		&i.Id,
		&i.Projectname,
		&i.CreatedOn,
		&i.CreatedBy,
		&i.Longitude,
		&i.Latitude,
		&i.Address,
		&i.Responsible,
		&i.Client,
		&i.CompanyId,
	)
	return i, err
}
