package db

import (
	"context"
	"testing"


	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)



func TestGetUser(t *testing.T) {
	user1 := generateRandomUserWithRole(t, util.ROLES[1])
	require.NotEmpty(t, user1)

	username := user1.Username

	user2, err := testStore.GetUser(context.Background(), username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

}


func TestCreateUserPersistence(t *testing.T) {
	generateRandomUserWithRole(t,util.ROLES[1])

}
