package db

import (
	"context"
	"fmt"

	"github.com/kamalbowselvam/chaintask/models"
	"github.com/lib/pq"
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
	Name      string  `json:"owner"`
	Budget    float64 `json:"balance"`
	CreatedBy string  `json:"created_by"`
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

const deleteAccount = `
DELETE FROM tasks WHERE id = $1
`

func (q *Queries) DeleteTask(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getTask = `
SELECT id, name, budget, created_on, created_by, updated_on, updated_by, done FROM tasks
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTask(ctx context.Context, id int64) (models.Task, error) {
	row := q.db.QueryRowContext(ctx, getTask, id)
	var t models.Task
	err := row.Scan(
		&t.Id,
		&t.Name,
		&t.Budget,
		&t.CreatedOn,
		&t.CreatedBy,
		&t.UpdatedOn,
		&t.UpdatedBy,
		&t.Done,
	)
	return t, err
}

const getTaskList = `
SELECT id, name, budget, created_on, created_by, updated_on, updated_by, done FROM tasks
WHERE id=any($1)
`

func (q *Queries) GetTaskList(ctx context.Context, ids []int64) ([]models.Task, error) {
	rows, err := q.db.QueryContext(ctx, getTaskList, pq.Array(ids))
	res := []models.Task{}
	if err != nil {
		return res, err
	}
	for rows.Next() {
		// FIXME Maybe that method could be extracted ?
		var t models.Task
		err = rows.Scan(
			&t.Id,
			&t.Name,
			&t.Budget,
			&t.CreatedOn,
			&t.CreatedBy,
			&t.UpdatedOn,
			&t.UpdatedBy,
			&t.Done,
		)
		res = append(res, t)
	}
	return res, err
}

const updateTask = `
 UPDATE tasks set name = $1, budget = $2, created_on = $3, created_by = $4, updated_on = $5, updated_by = $6, done = $7 where id = $8
`

func (q *Queries) UpdateTask(ctx context.Context, task models.Task) (models.Task, error) {
	// Create a helper function for preparing failure results.
	fail := func(err error) (models.Task, error) {
		return models.Task{}, fmt.Errorf("could not create Task: %v", err)
	}
	// Get a Tx for making transaction requests.
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	// Defer a rollback in case anything fails.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var id = task.Id
	_, err = tx.ExecContext(ctx, updateTask, task.Name, task.Budget, task.CreatedOn, task.CreatedBy, task.UpdatedOn, task.UpdatedBy, task.Done, id)
	if err != nil {
		return fail(err)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return q.GetTask(ctx, id)
}
