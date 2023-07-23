package db

import (
	"context"
	"testing"
	"time"

	"github.com/kamalbowselvam/chaintask/models"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func generateRandomTask(t *testing.T) models.Task {

	name := util.RandomName()
	arg := CreateTaskParams{
		Name:      name,
		Budget:    util.RandomBudget(),
		CreatedBy: name,
	}

	task, err := testQueries.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)
	require.Equal(t, arg.Name, task.Name)
	require.Equal(t, arg.Budget, task.Budget)
	require.Equal(t, arg.CreatedBy, task.CreatedBy)

	require.NotZero(t, task.Id)
	require.NotZero(t, task.CreatedOn)

	return task

}

func TestCreateTask(t *testing.T) {
	generateRandomTask(t)

}

func TestGetTask(t *testing.T) {
	task1 := generateRandomTask(t)

	require.NotEmpty(t, task1)

	task2, err := testQueries.GetTask(context.Background(), task1.Id)

	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task1.Name, task2.Name)
	require.Equal(t, task1.Budget, task2.Budget)
	require.Equal(t, task1.CreatedBy, task2.CreatedBy)
	require.WithinDuration(t, task1.CreatedOn, task2.CreatedOn, time.Second)

}

func TestDeleteTask(t *testing.T) {
	task1 := generateRandomTask(t)
	require.NotEmpty(t, task1)

	err := testQueries.DeleteTask(context.Background(), task1.Id)
	require.NoError(t, err)

}


func TestGetTaskList(t *testing.T) {
	task1 := generateRandomTask(t)
	task2 := generateRandomTask(t)
	task3 := generateRandomTask(t)
	taskList1, err := testQueries.GetTaskList(context.Background(), []int64{task1.Id, task2.Id})
	require.NoError(t, err)
	require.NotEmpty(t, taskList1)
	require.Equal(t, len(taskList1), 2)
	taskList2, err := testQueries.GetTaskList(context.Background(), []int64{task2.Id, task3.Id})
	require.NoError(t, err)
	require.NotEmpty(t, taskList2)
	require.Equal(t, len(taskList2), 2)
	taskList3, err := testQueries.GetTaskList(context.Background(), []int64{task2.Id + 1000, task3.Id + 1000})
	require.NoError(t, err)
	require.Empty(t, taskList3)

func TestUpdateTask(t *testing.T) {
	task1 := generateRandomTask(t)
	require.NotEmpty(t, task1)
	g := &task1
	g.Done = true
	g.Name = "test"
	require.Equal(t, task1.Done, true)
	task2, err := testQueries.UpdateTask(context.Background(), task1)
	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task2.Name, "test")
	require.Equal(t, task2.Done, true)

}
