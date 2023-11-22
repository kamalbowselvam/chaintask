package db

import (
	"context"
	"testing"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func generateRandomResource(t *testing.T) domain.Resource {
	arg := CreateResourceParams{
		ResourceName: util.RandomName(),
		Availed:      util.RandomBudget(),
		CreatedBy:    util.RandomString(10),
	}

	resource, err := testStore.CreateResource(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, resource)
	require.Equal(t, arg.ResourceName, resource.ResourceName)
	val1 ,_ := arg.Availed.Float64()
	val2, _ := arg.Availed.Float64()
	require.Equal(t, val1, val2)

	require.Equal(t, arg.CreatedBy, resource.CreatedBy)
	require.Equal(t, arg.CreatedBy, resource.UpdatedBy)

	require.NotZero(t, resource.Id)
	require.NotZero(t, resource.CreatedOn)
	require.NotZero(t, resource.UpdatedOn)

	return resource
}

func TestCreateResource(t *testing.T) {
	generateRandomResource(t)
}
