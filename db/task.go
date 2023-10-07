package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

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
  RETURNING id, taskname, budget, created_by, created_at, updated_by, updated_on, done, task_order, project_id;`

type CreateTaskParams struct {
	TaskName  string  `json:"taskname"`
	Budget    float64 `json:"budget"`
	CreatedBy string  `swaggerignore:"true"`
	TaskOrder int64   `json:"taskOrder"`
	ProjectId int64   `json:"projectId"`
}

type UpdateTaskParams struct {
	Id        int64     `json:"id"`
	TaskName  string    `json:"taskname"`
	Budget    float64   `json:"budget"`
	UpdatedOn time.Time `swaggerignore:"true"`
	UpdatedBy string    `swaggerignore:"true"`
	Done      bool      `json:"done"`
	TaskOrder int64     `json:"taskOrder"`
	ProjectId int64     `json:"ProjectId"`
	Version   int64     `json:"Version"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (domain.Task, error) {
	q.logger.Info("saving tasks")
	
	q.logger.Debug("Argument to Create task", zap.String("task_name",arg.TaskName),
	zap.Float64("budget",arg.Budget),
	zap.String("task_name",arg.CreatedBy),
	zap.Int64("task_name",arg.TaskOrder),
	zap.Int64("task_name",arg.ProjectId),
	)
	
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

const getTask = `SELECT id, taskname, budget, created_at, created_by, updated_on, updated_by, done, task_order, project_id FROM tasks
	WHERE id = $1 LIMIT 1
	`

type GetTaskParams struct {
	Id int64 `uri:"taskId" binding:"required,min=1"`
}

func (q *Queries) GetTask(ctx context.Context, id int64) (domain.Task, error) {
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
SELECT id, taskname, budget, created_at, created_by, updated_on, updated_by, done, task_order, project_id FROM tasks
WHERE id=any($1)
`

func (q *Queries) GetTaskList(ctx context.Context, ids []int64) ([]domain.Task, error) {
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
SELECT id, taskname, budget, created_at, created_by, updated_on, updated_by, done, task_order, project_id FROM tasks
WHERE project_id=$1 ORDER BY task_order ASC
`

func (q *Queries) GetTaskListByProject(ctx context.Context, project_id int64) ([]domain.Task, error) {
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

func (q *Queries) DeleteTask(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const deleteTaskFromProject = `DELETE FROM tasks WHERE project_id = $1`


func (q *Queries) DeleteTasksLinkedToProject(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTaskFromProject, id)
	return err
}

const updateTask = `
 UPDATE tasks set taskname = $1, budget = $2, updated_on = $3, updated_by = $4, done = $5, task_order=$6, project_id=$7, version=$9 + 1 where id = $8 and version = $9
`

func (q *Queries) UpdateTask(ctx context.Context, task UpdateTaskParams) (domain.Task, error) {
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
	row := tx.QueryRowContext(ctx, "SELECT id FROM tasks WHERE id=$1 limit 1", id)
	var oldId int64
	if err := row.Scan(&oldId); err != nil {
		return fail(err)
	}
	result, err := tx.ExecContext(ctx, updateTask, task.TaskName, task.Budget, task.UpdatedOn, task.UpdatedBy, task.Done, task.TaskOrder, task.ProjectId, id, task.Version)
	if err != nil {
		return fail(err)
	}
	affected, err  := result.RowsAffected() 
	if err != nil{
		return fail(err)
	}
	if affected != 1{
		return fail(errors.New("more than 1 row or 0 affected"));
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return q.GetTask(ctx, id)
}
