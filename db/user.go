package db

import (
	"context"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	"go.uber.org/zap"
)

const createUser = `INSERT INTO users (
	username, 
	hashed_password, 
	full_name, 
	email,
	user_role,
	company_id
) VALUES ( 
	$1, $2, $3, $4, UPPER($5), $6
) 
RETURNING username, hashed_password, full_name, email, created_at, user_role, company_id
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	UserRole       string `json:"user_role"`
	CompanyId      int64  `json:"company_id"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (domain.User, error) {

	logger_ := logger.FromCtx(ctx)
	logger_.Debug("Arguments to create user",
		zap.String("user_name", arg.Username),
		zap.String("hashed_password", arg.HashedPassword),
		zap.String("full_name", arg.FullName),
		zap.String("email", arg.Email),
		zap.String("user_role", arg.UserRole),
	    zap.Int64("company_id", arg.CompanyId))

	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.HashedPassword, arg.FullName, arg.Email, arg.UserRole, arg.CompanyId)
	var i domain.User

	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.CreatedAt,
		&i.UserRole,
		&i.CompanyId,
	)

	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, hashed_password, full_name, email, password_changed_at, created_at, user_role, company_id FROM users 
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (domain.User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i domain.User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.UserRole,
		&i.CompanyId,
	)
	return i, err
}

const deleteUser = `DELETE FROM users WHERE username = $1`

func (q *Queries) DeleteUser(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, username)
	return err
}
