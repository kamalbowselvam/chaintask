package repositories

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kamalbowselvam/chaintask/internal/core/domain"
	"github.com/kamalbowselvam/chaintask/internal/core/ports"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var testPersistenceStore *PersistenceSotrage
var testInMemoryStore *InMemoryStorage
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("Failed to load the config file")
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connet to db: ", err)
	}

	testPersistenceStore = NewPersistenceStorage(testDB)
	testInMemoryStore = NewInMemoryStorage()
	os.Exit(m.Run())

}

func generateRandomUser(t *testing.T, store ports.TaskRepository) domain.UserDetail {
	username := util.RandomName()
	hpassword, _ := util.HashPassword(util.RandomString(32))
	fname := util.RandomName()
	email := util.RandomEmail()
	user := domain.NewUser(username, hpassword, fname, email)
	userdetail, err := store.CreateUser(context.Background(), user)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return userdetail

}

func generateRandomTask(t *testing.T, store ports.TaskRepository) domain.Task {

	taskname := util.RandomName()
	budget := util.RandomBudget()
	userdetail  := generateRandomUser(t,store)
	task := domain.NewTask(taskname, budget, userdetail.Username)

	task, err := store.SaveTask(context.Background(), task)
	require.NoError(t, err)
	require.NotEmpty(t, task)
	require.Equal(t, taskname, task.TaskName)
	require.Equal(t, budget, task.Budget)
	require.Equal(t, userdetail.Username, task.CreatedBy)

	require.NotZero(t, task.Id)
	require.NotZero(t, task.CreatedOn)

	return task

}

func GetTaskHelper(t *testing.T, store ports.TaskRepository) {
	task1 := generateRandomTask(t, store)

	require.NotEmpty(t, task1)

	task2, err := store.GetTask(context.Background(), task1.Id)
	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task1.TaskName, task2.TaskName)
	require.Equal(t, task1.Budget, task2.Budget)
	require.Equal(t, task1.CreatedBy, task2.CreatedBy)
	require.WithinDuration(t, task1.CreatedOn, task2.CreatedOn, time.Second)

}

func GetTaskListHelper(t *testing.T, store ports.TaskRepository) {
	task1 := generateRandomTask(t, store)
	task2 := generateRandomTask(t, store)
	task3 := generateRandomTask(t, store)
	taskList1, err := store.GetTaskList(context.Background(), []int64{task1.Id, task2.Id})
	require.NoError(t, err)
	require.NotEmpty(t, taskList1)
	require.Equal(t, len(taskList1), 2)
	taskList2, err := store.GetTaskList(context.Background(), []int64{task2.Id, task3.Id})
	require.NoError(t, err)
	require.NotEmpty(t, taskList2)
	require.Equal(t, len(taskList2), 2)
	taskList3, err := store.GetTaskList(context.Background(), []int64{task2.Id + 1000, task3.Id + 1000})
	require.NoError(t, err)
	require.Empty(t, taskList3)
}

func DeleteTaskHelper(t *testing.T, store ports.TaskRepository) {
	task1 := generateRandomTask(t, store)
	require.NotEmpty(t, task1)

	err := store.DeleteTask(context.Background(), task1.Id)
	require.NoError(t, err)

}

func UpdateTaskHelper(t *testing.T, store ports.TaskRepository) {
	task1 := generateRandomTask(t, store)
	require.NotEmpty(t, task1)
	g := &task1
	g.Done = true
	g.TaskName = "test"
	require.Equal(t, task1.Done, true)
	task2, err := store.UpdateTask(context.Background(), task1)
	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task2.TaskName, "test")
	require.Equal(t, task2.Done, true)

}
