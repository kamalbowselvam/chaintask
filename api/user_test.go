package api

import (
	"testing"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)


func randomUser(t *testing.T) (user domain.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = domain.User{
		Username:       util.RandomName(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	}
	return
}
