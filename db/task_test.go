package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {

	arg := CreateTaskParams{
		Name:      "Kamal",
		Budget:    10000,
		CreatedBy: "Kamal",
	}

	task, err := testQueries.CreateTask(context.Background(), arg)
	if err != nil {
		log.Fatal(err.Error())
	}
	require.NoError(t, err)
	require.NotEmpty(t, task)
	require.Equal(t, arg.Name, task.Name)
	require.Equal(t, arg.Budget, task.Budget)
	require.Equal(t, arg.CreatedBy, task.CreatedBy)

	require.NotZero(t, task.Id)
	require.NotZero(t, task.CreatedOn)

}
