package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kamalbowselvam/chaintask/internal/core/domain"
	"github.com/lib/pq"
)


type PersistenceSotrage struct {
	db *sql.DB
}

func NewPersistenceStorage(db *sql.DB) *PersistenceSotrage {
	return &PersistenceSotrage{
		db: db,
	}
}

const createTask = `INSERT INTO tasks (
	name,
	budget,
	created_by,
	updated_by 
  ) VALUES (
	$1, $2, $3, $4
  )
  RETURNING id, name, budget, created_by, created_on, updated_by, updated_on, done;`
  
  
func (q *PersistenceSotrage) SaveTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	  row := q.db.QueryRowContext(ctx, createTask, task.Name, task.Budget, task.CreatedBy, task.UpdatedBy)
	  var i domain.Task
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


const getTask = `SELECT id, name, budget, created_on, created_by, updated_on, updated_by, done FROM tasks
	WHERE id = $1 LIMIT 1
	`
	
func (q *PersistenceSotrage) GetTask(ctx context.Context, id int64) (domain.Task, error) {
		row := q.db.QueryRowContext(ctx, getTask, id)
		var t domain.Task
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

func (q *PersistenceSotrage) GetTaskList(ctx context.Context, ids []int64) ([]domain.Task, error) {
	rows, err := q.db.QueryContext(ctx, getTaskList, pq.Array(ids))
	res := []domain.Task{}
	if err != nil {
		return res, err
	}
	for rows.Next() {
		// FIXME Maybe that method could be extracted ?
		var t domain.Task
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


const deleteAccount = `DELETE FROM tasks WHERE id = $1`

func (q *PersistenceSotrage) DeleteTask(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}



const updateTask = `
 UPDATE tasks set name = $1, budget = $2, created_on = $3, created_by = $4, updated_on = $5, updated_by = $6, done = $7 where id = $8
`

func (q *PersistenceSotrage) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	// Create a helper function for preparing failure results.
	fail := func(err error) (domain.Task, error) {
		return domain.Task{}, fmt.Errorf("could not create Task: %v", err)
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
