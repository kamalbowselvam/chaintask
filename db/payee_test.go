package db

import (
	"context"
	"testing"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func generateRandomPayee(t *testing.T) domain.Payee {
	arg := CreatePayeeParams{
		PayeeName: util.RandomName(),
		Address: util.RandomAddress(),
		CreatedBy: util.RandomString(10),
	}

	payee, err := testStore.CreatePayee(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payee)
	require.Equal(t, arg.PayeeName, payee.PayeeName)
	require.Equal(t, arg.Address, payee.Address)
	require.Equal(t, arg.CreatedBy, payee.CreatedBy)
	require.Equal(t, arg.CreatedBy, payee.UpdatedBy)

	require.NotZero(t, payee.Id)
	require.NotZero(t, payee.CreatedOn)
	require.NotZero(t, payee.UpdatedOn)

	return payee
}

func TestCreatePayee(t *testing.T) {
	generateRandomPayee(t)
}