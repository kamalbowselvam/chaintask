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
	role_id 
) VALUES ( 
	$1, $2, $3, $4, (select id from roles where userRole=UPPER($5))
) 
RETURNING username, hashed_password, full_name, email, created_at, role_id
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (domain.User, error) {

	q.logger.Debug("Arguments to create user", 
				zap.String("user_name", arg.Username),
				zap.String("hashed_password", arg.HashedPassword),
				zap.String("full_name", arg.FullName),
				zap.String("email", arg.Email), 
				zap.String("role", arg.Role),)



	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.HashedPassword, arg.FullName, arg.Email, arg.Role)
	var i domain.User

	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.CreatedAt,
		&i.Role,
	)

	return i, err
}




const getUser = `-- name: GetUser :one
SELECT username, hashed_password, full_name, email, password_changed_at, created_at, role_id as role FROM users left join roles on role_id = roles.id 
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
		&i.Role,
	)
	return i, err
}