package domain

import "time"

type User struct {
	Id                int64
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	UserRole          string    `json:"user_role"`
	CompanyId         int64     `json:"company_id"`
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
