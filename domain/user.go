package domain

import "time"

type User struct {
	Id                int64
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"CreatedAt"`
	PasswordChangedAt time.Time `json:"PasswordChangedAt"`
	UserRole          string    `json:"user_role"`
}

func NewUser(username string, hpassord string, fname string, email string, role string) User {

	return User{
		Username:       username,
		HashedPassword: hpassord,
		FullName:       fname,
		Email:          email,
		UserRole:       role,
	}
}
