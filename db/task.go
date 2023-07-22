package db

import (
	"context"
	"fmt"
	"reflect"

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

const updateTask = `
 UPDATE tasks set $1 = $2 where id = $3
`

func (q *Queries) UpdateTask(ctx context.Context, task models.Task) (models.Task, error) {
	// Create a helper function for preparing failure results.
	fail := func(err error) (models.Task, error) {
		return models.Task{}, fmt.Errorf("Could not create Task: %v", err)
	}
	// Get a Tx for making transaction requests.
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()
	var id = task.Id
	reflectedTask := reflect.ValueOf(task)
	for i := 0; i < reflectedTask.NumField(); i++ {
		if !reflect.Indirect(reflectedTask).FieldByName(reflectedTask.Type().Field(i).Name).IsNil() {
			tx.ExecContext(ctx, updateTask, reflectedTask.Type().Field(i).Name, reflect.Indirect(reflectedTask).FieldByName(reflectedTask.Type().Field(i).Name), id)
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return q.GetTask(ctx, id)
}
