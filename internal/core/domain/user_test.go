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