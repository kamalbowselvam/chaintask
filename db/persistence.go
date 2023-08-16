package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kamalbowselvam/chaintask/domain"
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


const createUser = `INSERT INTO users (
	username, 
	hashed_password, 
	full_name, 
	email,
	role_id 
) VALUES ( 
	$1, $2, $3, $4, (select id from roles where userRole=$5)
) 
RETURNING username, hashed_password, full_name, email, created_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
}

func (q *PersistenceSotrage) CreateUser(ctx context.Context, arg CreateUserParams) (domain.User, error){
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.HashedPassword, arg.FullName, arg.Email, arg.Role)
	var i domain.User

	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.CreatedAt,
	)

	return i, err
}

const createTask = `INSERT INTO tasks (
	taskname,
	budget,
	created_by,
	updated_by 
  ) VALUES (
	$1, $2, $3, $4
  )
  RETURNING id, taskname, budget, created_by, created_on, updated_by, updated_on, done;`
  

type CreateTaskParams struct {
	TaskName  string  `json:"taskname"`
	Budget    float64 `json:"budget"`
	CreatedBy string  `json:"createdBy"`
}

  
func (q *PersistenceSotrage) CreateTask(ctx context.Context, arg CreateTaskParams) (domain.Task, error) {
	  row := q.db.QueryRowContext(ctx, createTask, arg.TaskName, arg.Budget, arg.CreatedBy, arg.CreatedBy)
	  var i domain.Task
	  err := row.Scan(
		  &i.Id,
		  &i.TaskName,
		  &i.Budget,
		  &i.CreatedBy,
		  &i.CreatedOn,
		  &i.UpdatedBy,
		  &i.UpdatedOn,
		  &i.Done,
	  )
	  return i, err

	}


const getTask = `SELECT id, taskname, budget, created_on, created_by, updated_on, updated_by, done FROM tasks
	WHERE id = $1 LIMIT 1
	`
	
func (q *PersistenceSotrage) GetTask(ctx context.Context, id int64) (domain.Task, error) {
		row := q.db.QueryRowContext(ctx, getTask, id)
		var t domain.Task
		err := row.Scan(
			&t.Id,
			&t.TaskName,
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
SELECT id, taskname, budget, created_on, created_by, updated_on, updated_by, done FROM tasks
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
			&t.TaskName,
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
 UPDATE tasks set taskname = $1, budget = $2, created_on = $3, created_by = $4, updated_on = $5, updated_by = $6, done = $7 where id = $8
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
	_, err = tx.ExecContext(ctx, updateTask, task.TaskName, task.Budget, task.CreatedOn, task.CreatedBy, task.UpdatedOn, task.UpdatedBy, task.Done, id)
	if err != nil {
		return fail(err)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return q.GetTask(ctx, id)
}


const getUser = `-- name: GetUser :one
SELECT username, hashed_password, full_name, email, password_changed_at, created_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *PersistenceSotrage) GetUser(ctx context.Context, username string) (domain.User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i domain.User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}