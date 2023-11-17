package db

import (
	"context"

	"github.com/kamalbowselvam/chaintask/domain"
	"go.uber.org/zap"
)

const createUser = `INSERT INTO users (
	username, 
	hashed_password, 
	full_name, 
	email,
	user_role
) VALUES ( 
	$1, $2, $3, $4, UPPER($5)
) 
RETURNING username, hashed_password, full_name, email, created_at, user_role
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	UserRole       string `json:"user_role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (domain.User, error) {

	q.logger.Debug("Arguments to create user",
		zap.String("user_name", arg.Username),
		zap.String("hashed_password", arg.HashedPassword),
		zap.String("full_name", arg.FullName),
		zap.String("email", arg.Email),
		zap.String("user_role", arg.UserRole))

	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.HashedPassword, arg.FullName, arg.Email, arg.UserRole)
	var i domain.User

	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.CreatedAt,
		&i.UserRole,
	)

	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, hashed_password, full_name, email, password_changed_at, created_at, user_role FROM users 
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
	)
	return i, err
}

const deleteUser = `DELETE FROM users WHERE username = $1`

func (q *Queries) DeleteUser(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, username)
	return err
}
