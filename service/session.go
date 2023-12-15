package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
)

func (srv *service) CreateSession(ctx context.Context, arg db.CreateSessionParams) (domain.Session, error) {
	session, err := srv.globalRepository.CreateSession(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), arg)
	return session, err

}

func (srv *service) GetSession(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	session, err := srv.globalRepository.GetSession(logger.WithCtx(context.Background(), logger.FromCtx(ctx)), id)
	return session, err

}
