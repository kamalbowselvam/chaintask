package db

import (
	"context"

	"github.com/kamalbowselvam/chaintask/domain"
	"go.uber.org/zap"
)

const createCompany = `
  INSERT into company (companyname, address, created_by, updated_by) values ($1, $2, $3, $3)
  RETURNING id, companyname, address, created_at, created_by, updated_on, updated_by
`

type CreateCompanyParams struct {
	CompanyName string  `json:"companyname"`
	Address     string  `json:"address"`
	CreatedBy   string
}

func (q *Queries) CreateCompany(ctx context.Context, arg CreateCompanyParams) (domain.Company, error){
	q.logger.Info("saving comapny")
	
	q.logger.Debug("Argument to Create a company", zap.String("company_name",arg.CompanyName),
	zap.String("address",arg.Address),
	)
	
	row := q.db.QueryRowContext(ctx, createCompany, arg.CompanyName, arg.Address, arg.CreatedBy)
	var i domain.Company
	err := row.Scan(
		&i.Id,
		&i.CompanyName,
		&i.Address,
		&i.CreatedOn,
		&i.CreatedBy,
		&i.UpdatedOn,
		&i.UpdatedBy,
	)
	return i, err
}