package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"


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
RETURNING username, hashed_password, full_name, email, created_at, role_id
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
}

func (q *PersistenceSotrage) CreateUser(ctx context.Context, arg CreateUserParams) (domain.User, error) {
	log.Println(arg)
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.HashedPassword, arg.FullName, arg.Email, arg.Role)
	var i domain.User

	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.CreatedAt,
		&i.Role,
	)

	return i, err
}

const createTask = `INSERT INTO tasks (
	taskname,
	budget,
	created_by,
	updated_by,
	task_order,
	project_id 
  ) VALUES (
	$1, $2, $3, $4, $5, $6
  )
  RETURNING id, taskname, budget, created_by, created_on, updated_by, updated_on, done, task_order, project_id;`

type CreateTaskParams struct {
	TaskName  string  `json:"taskname"`
	Budget    float64 `json:"budget"`
	CreatedBy string  `json:"createdBy"`
	TaskOrder int64   `json:"taskOrder"`
	ProjectId int64   `json:"projectId"`
}

func (q *PersistenceSotrage) CreateTask(ctx context.Context, arg CreateTaskParams) (domain.Task, error) {
	log.Println("saving tasks")
	log.Println(arg)
	row := q.db.QueryRowContext(ctx, createTask, arg.TaskName, arg.Budget, arg.CreatedBy, arg.CreatedBy, arg.TaskOrder, arg.ProjectId)
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
		&i.TaskOrder,
		&i.ProjectId,
	)
	return i, err

}

const getTask = `SELECT id, taskname, budget, created_on, created_by, updated_on, updated_by, done, task_order, project_id FROM tasks
	WHERE id = $1 LIMIT 1
	`

type GetTaskParams struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

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
		&t.TaskOrder,
		&t.ProjectId,
	)
	return t, err
}

const getTaskList = `
SELECT id, taskname, budget, created_on, created_by, updated_on, updated_by, done, task_order, project_id FROM tasks
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
			&t.TaskOrder,
			&t.ProjectId,
		)
		res = append(res, t)
	}
	return res, err
}


const getTaskListByProject = `
SELECT id, taskname, budget, created_on, created_by, updated_on, updated_by, done, task_order, project_id FROM tasks
WHERE project_id=$1 ORDER BY task_order ASC
`

func (q *PersistenceSotrage) GetTaskListByProject(ctx context.Context, project_id int64) ([]domain.Task, error) {
	rows, err := q.db.QueryContext(ctx, getTaskListByProject, project_id)
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
			&t.TaskOrder,
			&t.ProjectId,
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
 UPDATE tasks set taskname = $1, budget = $2, created_on = $3, created_by = $4, updated_on = $5, updated_by = $6, done = $7, task_order=$8, project_id=$9 where id = $10
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
	_, err = tx.ExecContext(ctx, updateTask, task.TaskName, task.Budget, task.CreatedOn, task.CreatedBy, task.UpdatedOn, task.UpdatedBy, task.Done, task.TaskOrder, task.ProjectId, id)
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
SELECT username, hashed_password, full_name, email, password_changed_at, created_at, role_id as role FROM users left join roles on role_id = roles.id 
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
		&i.Role,
	)
	return i, err
}

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

func (q *PersistenceSotrage) CreateProject(ctx context.Context, arg CreateProjectParam) (domain.Project, error) {
	log.Println("saving projects")
	log.Println(arg)
	row := q.db.QueryRowContext(ctx, createProject, arg.ProjectName, arg.CreatedBy, Point{arg.Location[0], arg.Location[1]}, arg.Address, arg.Responsible, arg.Client)
	var i domain.Project
	var p Point;
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
