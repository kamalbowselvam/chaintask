package db

import (
	"context"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const createTaskPayment = `INSERT INTO tasks_payments (
	task_id,
	amount,
	payee,
	created_by,
	updated_by
  ) VALUES (
	$1, $2, $3, $4, $4
  )
  RETURNING id, task_id, amount, payee, created_by, created_at, updated_by, updated_at`

type CreateTaskPaymentParams struct {
	TaskId    int64           `json:"taskIf" binding:"required,number"`
	Amount    decimal.Decimal `json:"amount" binding:"required,number"`
	Source    string          `json:"source" binding:"required"`
	CreatedBy string
}

func (q *Queries) PayForATask(ctx context.Context, arg CreateTaskPaymentParams) (domain.TaskPayment, error) {

	logger_ := logger.FromCtx(ctx)
	logger_.Debug("Creating a task",
		zap.String("package", "db"),
		zap.String("function", "CreateTaskPayment"),
		zap.Any("param", arg),
	)

	row := q.db.QueryRowContext(ctx, createTaskPayment, arg.TaskId, arg.Amount, arg.Source)
	var i domain.TaskPayment
	err := row.Scan(
		&i.Id,
		&i.TaskId,
		&i.Amount,
		&i.Source,
		&i.CreatedBy,
		&i.CreatedOn,
		&i.UpdatedBy,
		&i.UpdatedOn,
	)
	return i, err

}

const getPaymentsByTasks = `
SELECT id, task_id, amount, payee FROM tasks_payments WHERE task_id=$1 ORDER BY amount ASC`

func (q *Queries) getTaksPaymentsByTask(ctx context.Context, task_id int64) ([]domain.TaskPayment, error) {
	rows, err := q.db.QueryContext(ctx, getPaymentsByTasks, task_id)
	res := []domain.TaskPayment{}
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var t domain.TaskPayment
		err = rows.Scan(
			&t.Id,
			&t.TaskId,
			&t.Amount,
			&t.Source,
		)
		res = append(res, t)
	}
	return res, err
}
