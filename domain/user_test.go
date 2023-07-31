package domain

import (
	"testing"

	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	user := User{
		Username:       util.RandomName(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	}

	return user
}

func TestNewUser(t *testing.T) {

	hashedPassword, err := util.HashPassword(util.RandomString(6))

	require.NoError(t,err)

	username := util.RandomName()
	hpassword := hashedPassword
	fullname := util.RandomName()
	email := util.RandomEmail()

	user := NewUser(username, hpassword, fullname, email)

	require.NotEmpty(t, user)
	require.Equal(t, username, user.Username)
	require.Equal(t, hpassword, user.HashedPassword)
	require.Equal(t, fullname, user.FullName)
	require.Equal(t, email, user.Email)

}
