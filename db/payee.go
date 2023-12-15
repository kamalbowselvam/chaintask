package db

import (
	"context"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	"go.uber.org/zap"
)

const createPayee = `
  INSERT into payees (payee_name, address, created_by, updated_by) values ($1, $2, $3, $3)
  RETURNING id, payee_name, address, created_at, created_by, updated_on, updated_by
`

type CreatePayeeParams struct {
	PayeeName string `json:"payee_name"`
	Address   string `json:"address"`
	CreatedBy string
}

func (q *Queries) CreatePayee(ctx context.Context, arg CreatePayeeParams) (domain.Payee, error) {

	logger_ := logger.FromCtx(ctx)
	logger_.Debug("Argument to Create a resource", zap.String("company_name", arg.PayeeName),
		zap.String("address", arg.Address),
	)

	row := q.db.QueryRowContext(ctx, createPayee, arg.PayeeName, arg.Address, arg.CreatedBy)
	var i domain.Payee
	err := row.Scan(
		&i.Id,
		&i.PayeeName,
		&i.Address,
		&i.CreatedOn,
		&i.CreatedBy,
		&i.UpdatedOn,
		&i.UpdatedBy,
	)
	return i, err
}
