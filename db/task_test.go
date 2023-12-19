package db

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"log"
	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func generateRandomTask(t *testing.T) domain.Task {

	project := generateRandomProject(t)
	taskOrder := util.RandomInt(1, 100)
	arg := CreateTaskParams{
		TaskName:  util.RandomName(),
		Budget:    util.RandomBudget(),
		CreatedBy: project.Client,
		TaskOrder: taskOrder,
		ProjectId: project.Id,
	}

	task, err := testStore.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)
	require.Equal(t, arg.TaskName, task.TaskName)
	val1 ,_ := arg.Budget.Float64()
	val2, _ := arg.Budget.Float64()
	require.Equal(t, val1, val2)
	require.Equal(t, arg.CreatedBy, task.CreatedBy)
	require.Equal(t, arg.TaskOrder, task.TaskOrder)
	require.Equal(t, arg.ProjectId, task.ProjectId)
	require.Equal(t, task.CompanyId, project.CompanyId)

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
	val1 ,_ := task1.Budget.Float64()
	val2, _ := task2.Budget.Float64()
	require.Equal(t, val1, val2)
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
	version := int64(0)
	rating := int64(0)
	update := UpdateTaskParams{}
	update.Budget = task1.Budget
	update.Done = task1.Done
	update.ProjectId = &task1.ProjectId
	update.TaskOrder = task1.TaskOrder
	update.Version = &version
	update.Id = task1.Id
	update.UpdatedBy = task1.CreatedBy;
	update.UpdatedOn = time.Now()
	update.TaskName = "test"
	update.Done = true
	update.Rating = &rating
	task2, err := testStore.UpdateTask(context.Background(), update)
	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task2.TaskName, "test")
	require.Equal(t, task2.Done, true)

}

// Taken from https://hackernoon.com/comparing-optimistic-and-pessimistic-locking-with-go-and-postgresql
func TestOptimistic(t *testing.T){
	ctx := context.Background()
	task := generateRandomTask(t)
	user := generateRandomWorksManager(t)

	start := time.Now()
	if err := testSaveBulkData(ctx, task, user.Username); (err != nil){
		t.Errorf("Deposit() error = %v", err)
	}
	log.Println("Execution time", time.Since(start))
	
}

func testSaveBulkData(ctx context.Context, task domain.Task, userName string) error {
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go save(ctx, strconv.Itoa((i%100)+1), task, userName, sem, &wg)
	}
	wg.Wait()
	return nil
}

func save(ctx context.Context, taskName string, task domain.Task, userName string, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	sem <- struct{}{}
	// Logging here seems important, otherwise the go routine go in timeout. There is something shaddy to explore a bit further
	log.Println("Trying to update tasks")
	_, err := testStore.UpdateTask(ctx, UpdateTaskParams{TaskName: taskName, Id: task.Id, TaskOrder: task.TaskOrder, ProjectId: &task.ProjectId, UpdatedBy: userName})
	if err != nil {
		log.Printf("Error %s", err.Error())
	}
	<-sem
}