package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"go.uber.org/zap"
)

func (srv *service) CreateUser(ctx context.Context, arg db.CreateUserParams) (domain.User, error) {

	user, err := srv.globalRepository.CreateUser(context.Background(), arg)
	if err == nil {
		err = srv.policiesRepository.CreateUserPolicies(arg.Username, arg.UserRole, arg.CompanyId)
		if err != nil {
			err = srv.DeleteUser(ctx, arg.Username)
			if err != nil {
				srv.logger.Fatal("could not delete user ", zap.String("usernmae", arg.Username))
			}
			return domain.User{}, err
		}
	}
	return user, err

}

func (srv *service) GetUser(ctx context.Context, username string) (domain.User, error) {
	return srv.globalRepository.GetUser(context.Background(), username)
}

func (srv *service) DeleteUser(ctx context.Context, username string) error {
	return srv.globalRepository.DeleteUser(context.Background(), username)
}