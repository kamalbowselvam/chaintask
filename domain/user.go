package domain

import "time"

type User struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `josn:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"CreatedAt"`
	PasswordChangedAt time.Time `json:"PasswordChangedAt"`
}

func NewUser(username string, hpassord string, fname string, email string) User {

	return User{
		Username:       username,
		HashedPassword: hpassord,
		FullName:       fname,
		Email:          email,
	}
}
