package db

import (
	"context"
	"testing"
	"time"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func generateRandomTask(t *testing.T) domain.Task {

	project := generateRandomProject(t)
	arg := CreateTaskParams{
		TaskName:  util.RandomName(),
		Budget:    util.RandomBudget(),
		CreatedBy: project.Client,
		TaskOrder: util.RandomInt(0, 100),
		ProjectId: project.Id,
	}

	task, err := testStore.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)
	require.Equal(t, arg.TaskName, task.TaskName)
	require.Equal(t, arg.Budget, task.Budget)
	require.Equal(t, arg.CreatedBy, task.CreatedBy)
	require.Equal(t, arg.TaskOrder, task.TaskOrder)
	require.Equal(t, arg.ProjectId, task.ProjectId)

	require.NotZero(t, task.Id)
	require.NotZero(t, task.CreatedOn)

	return task

}

func TestGetTask(t *testing.T) {
	task1 := generateRandomTask(t)

	require.NotEmpty(t, task1)

	task2, err := testStore.GetTask(context.Background(), task1.Id)
	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task1.TaskName, task2.TaskName)
	require.Equal(t, task1.Budget, task2.Budget)
	require.Equal(t, task1.CreatedBy, task2.CreatedBy)
	require.Equal(t, task1.TaskOrder, task2.TaskOrder)
	require.Equal(t, task1.ProjectId, task2.ProjectId)
	require.WithinDuration(t, task1.CreatedOn, task2.CreatedOn, time.Second)

}

func TestCreateTaskPersistence(t *testing.T) {
	generateRandomTask(t)
}

func TestGetTaskList(t *testing.T) {
	task1 := generateRandomTask(t)
	task2 := generateRandomTask(t)
	task3 := generateRandomTask(t)

	taskList1, err := testStore.GetTaskList(context.Background(), []int64{task1.Id, task2.Id})
	require.NoError(t, err)
	require.NotEmpty(t, taskList1)
	require.Equal(t, len(taskList1), 2)
	taskList2, err := testStore.GetTaskList(context.Background(), []int64{task2.Id, task3.Id})
	require.NoError(t, err)
	require.NotEmpty(t, taskList2)
	require.Equal(t, len(taskList2), 2)
	taskList3, err := testStore.GetTaskList(context.Background(), []int64{task2.Id + 1000, task3.Id + 1000})
	require.NoError(t, err)
	require.Empty(t, taskList3)
}

func TestDeleteTask(t *testing.T) {
	task1 := generateRandomTask(t)
	require.NotEmpty(t, task1)

	err := testStore.DeleteTask(context.Background(), task1.Id)
	require.NoError(t, err)

}

func TestUpdateTaskHelper(t *testing.T) {
	task1 := generateRandomTask(t)
	require.NotEmpty(t, task1)
	update := UpdateTaskParams{}
	update.Budget = task1.Budget
	update.Done = task1.Done
	update.ProjectId = task1.ProjectId
	update.TaskOrder = task1.TaskOrder
	update.Version = 0
	update.Id = task1.Id
	update.UpdatedBy = task1.CreatedBy;
	update.UpdatedOn = time.Now()
	update.TaskName = "test"
	update.Done = true
	task2, err := testStore.UpdateTask(context.Background(), update)
	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task2.TaskName, "test")
	require.Equal(t, task2.Done, true)

}
