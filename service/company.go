package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	"go.uber.org/zap"
)

func (srv *service) CreateCompany(ctx context.Context, arg db.CreateCompanyParams) (domain.Company, error) {
	logger_ := logger.FromCtx(ctx)
	company, err := srv.globalRepository.CreateCompany(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), arg)
	if err != nil {
		logger_.Error("Could not save the task in repository", zap.Error(err))
		return company, err
	}
	return company, nil
}