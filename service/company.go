package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"go.uber.org/zap"
)

func (srv *service) CreateCompany(ctx context.Context, arg db.CreateCompanyParams) (domain.Company, error){
	company, err := srv.globalRepository.CreateCompany(context.Background(), arg)
	if err != nil {
		srv.logger.Fatal("Could not save the task in repository", zap.Error(err))
		return company, err
	}
	return company, nil
}