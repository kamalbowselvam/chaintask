package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)

func (srv *service) CreateSession(ctx context.Context, arg db.CreateSessionParams) (domain.Session, error) {
	session, err := srv.globalRepository.CreateSession(context.Background(), arg)
	return session, err

}

func (srv *service) GetSession(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	session, err := srv.globalRepository.GetSession(context.Background(), id)
	return session, err

}
