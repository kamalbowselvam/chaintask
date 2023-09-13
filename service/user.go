package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)

func (srv *service) CreateUser(ctx context.Context, arg db.CreateUserParams) (domain.User, error) {

	user, err := srv.globalRepository.CreateUser(context.Background(), arg)
	if err == nil {
		srv.policiesRepository.CreateUserPolicies(arg.Username, arg.Role)
		// FIXME check for error here and try to delete the user
	}
	return user, err

}

func (srv *service) GetUser(ctx context.Context, username string) (domain.User, error) {
	user, err := srv.globalRepository.GetUser(context.Background(), username)

	return user, err
}
