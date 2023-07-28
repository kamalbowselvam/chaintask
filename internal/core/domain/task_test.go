package domain

import (
	"testing"
	"time"

	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)


func createRandomTask(t *testing.T) Task {
	name := util.RandomName()
	budget := util.RandomBudget()
	createdBy := util.RandomName()

	task := NewTask(name,budget,createdBy)
	require.NotEmpty(t,task)
	require.Equal(t, task.TaskName, name)
	require.Equal(t, task.Budget, budget)
	require.Equal(t, task.CreatedBy, createdBy)
	require.Equal(t, task.UpdatedBy, createdBy)

	return task
}



func TestNewTask(t *testing.T){

	currentTime := time.Now()
	task := createRandomTask(t)
	require.IsType(t, task.CreatedOn, time.Now())
	require.IsType(t, task.UpdatedOn,time.Now())
	require.WithinDuration(t, task.CreatedOn,currentTime, time.Second)
	require.WithinDuration(t, task.UpdatedOn,currentTime, time.Second)
	require.Equal(t, task.Done, false)
}



func TestGetTaskName(t *testing.T){

	task := createRandomTask(t)
	name := task.GetTaskName()
	require.Equal(t,task.TaskName, name)

}

func TestGetTaskBudget(t *testing.T){

	task := createRandomTask(t)
	budget := task.GetBudget()
	require.Equal(t,task.Budget, budget)

}




func TestIsTaskDone(t *testing.T){

	task := createRandomTask(t)
	require.Equal(t, task.Done, false)
	task.SetTaskDone(true)
	require.Equal(t,task.Done, true)

}
