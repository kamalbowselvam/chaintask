package db

import (
	"context"
	"github.com/kamalbowselvam/chaintask/models"
)


const createTask = `
INSERT INTO tasks (
  name,
  budget,
  created_by,
  updated_by 
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, budget, created_by, created_on, updated_by, updated_on, done;`

type CreateTaskParams struct {
	Name    string `json:"owner"`
	Budget  float64 `json:"balance"`
	CreatedBy string `json:"created_by"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (models.Task, error) {
	row := q.db.QueryRowContext(ctx, createTask, arg.Name, arg.Budget, arg.CreatedBy, arg.CreatedBy)
	var i models.Task
	err := row.Scan(
		&i.Id,
		&i.Name,
		&i.Budget,
		&i.CreatedBy,
		&i.CreatedOn,
		&i.UpdatedBy,
		&i.UpdatedOn,
		&i.Done,
	)
	return i, err
}
