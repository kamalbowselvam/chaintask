package db

import (
	"context"
	"github.com/kamalbowselvam/chaintask/domain"
	"go.uber.org/zap"
)

const createResource = `
  INSERT into resources (resource_name, availed, created_by, updated_by) values ($1, $2, $3, $3)
  RETURNING id, resource_name, availed, current, created_at, created_by, updated_on, updated_by
`

type CreateResourceParams struct {
	ResourceName string  `json:"resource_name"`
	Availed 	float64  `json:"availed"`
	CreatedBy   string
}

func (q *Queries) CreateResource(ctx context.Context, arg CreateResourceParams) (domain.Resource, error){
	
	q.logger.Debug("Argument to Create a resource", zap.String("company_name",arg.ResourceName),
	zap.Float64("availed",arg.Availed),
	)
	
	row := q.db.QueryRowContext(ctx, createResource, arg.ResourceName, arg.Availed, arg.CreatedBy)
	var i domain.Resource
	err := row.Scan(
		&i.Id,
		&i.ResourceName,
		&i.Availed,
		&i.Current,
		&i.CreatedOn,
		&i.CreatedBy,
		&i.UpdatedOn,
		&i.UpdatedBy,
	)
	return i, err
}