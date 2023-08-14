package domain

import "time"

type User struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `josn:"full_name"`
	Email          string `json:"email"`
}

type UserDetail struct {
	Username  string    `json:"username"`
	FullName  string    `josn:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_by"`
}

func NewUser(username string, hpassord string, fname string, email string) User {

	return User{
		Username:       username,
		HashedPassword: hpassord,
		FullName:       fname,
		Email:          email,
	}
}

func NewUserDetail(username string, fname string, email string, createdAt time.Time) UserDetail {

	return UserDetail{
		Username:  username,
		FullName:  fname,
		Email:     email,
		CreatedAt: createdAt,
	}
}
