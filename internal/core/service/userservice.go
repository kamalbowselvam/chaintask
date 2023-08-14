package service

import (
	"context"
	"log"

	"github.com/kamalbowselvam/chaintask/internal/core/domain"
	"github.com/kamalbowselvam/chaintask/internal/core/ports"
)

type userService struct {
	userRepository ports.UserRepository
}

func NewUserService(userRepository ports.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (srv *userService) CreateUser(username string, hpassord string, fname string, email string) (domain.UserDetail, error) {

	user := domain.NewUser(username, hpassord, fname, email)
	userdetail, err := srv.userRepository.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatal("Could not save the user in repository", err.Error())

	}

	return userdetail, err

}

func (srv *userService) GetUser(username string) (domain.User, error) {
	user, err := srv.userRepository.GetUser(context.Background(), username)
	if err != nil {
		log.Fatalf("Could not find user with %s", username)
	}
	return user, err
}
