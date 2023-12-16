package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	"go.uber.org/zap"
)

func (srv *service) CreateUser(ctx context.Context, arg db.CreateUserParams) (domain.User, error) {

	logger_ := logger.FromCtx(ctx)
	user, err := srv.globalRepository.CreateUser(logger.WithCtx(context.Background(), logger_), arg)
	if err == nil {
		err = srv.policiesRepository.CreateUserPolicies(arg.Username, arg.UserRole, arg.CompanyId)
		if err != nil {
			err = srv.DeleteUser(logger.WithCtx(context.Background(), logger_), arg.Username)
			if err != nil {
				logger.Fatal("could not delete user ", zap.String("usernmae", arg.Username))
			}
			return domain.User{}, err
		}
	}
	return user, err

}

func (srv *service) GetUser(ctx context.Context, username string) (domain.User, error) {
	return srv.globalRepository.GetUser(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), username)
}

func (srv *service) DeleteUser(ctx context.Context, username string) error {
	return srv.globalRepository.DeleteUser(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), username)
}