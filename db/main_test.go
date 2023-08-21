package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var testPersistenceStore *PersistenceSotrage
var testInMemoryStore *InMemoryStorage
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../")

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

func generateRandomUserWithRole(t *testing.T, store GlobalRepository, role string) domain.User {

	hpassword, _ := util.HashPassword(util.RandomString(32))
	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: hpassword,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
		Role:           role,
	}
	user, err := store.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user

}

func generateRandomUser(t *testing.T, store GlobalRepository) domain.User {
	return generateRandomUserWithRole(t, store, util.RandomRoleString())
}

func generateRandomWorksManager(t *testing.T, store GlobalRepository) domain.User {
	return generateRandomUserWithRole(t, store, util.ROLES[2])
}

func generateRandomClient(t *testing.T, store GlobalRepository) domain.User {
	return generateRandomUserWithRole(t, store, util.ROLES[1])
}

func generateRandomLocation() domain.Location {
	return domain.Location{
		util.RandomLatitude(), util.RandomLongitude()}
}

func generateRandomProject(t *testing.T, store GlobalRepository) domain.Project {
	resp := generateRandomWorksManager(t, store)
	client := generateRandomClient(t, store)
	arg := CreateProjectParam{
		ProjectName: util.RandomName(),
		CreatedOn:   time.Now(),
		CreatedBy:   resp.Username,
		Client:      client.Username,
		Responsible: resp.Username,
		Address:     util.RandomAddress(),
		Location:    generateRandomLocation(),
	}

	project, err := store.CreateProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, project)
	require.Equal(t, arg.ProjectName, project.Projectname)
	require.WithinDuration(t, arg.CreatedOn, project.CreatedOn, time.Second)
	require.Equal(t, arg.CreatedBy, project.CreatedBy)
	require.Equal(t, arg.Client, project.Client)
	require.Equal(t, arg.Responsible, project.Responsible)
	require.NotZero(t, project.Id)
	require.NotZero(t, project.CreatedOn)
	return project
}

func generateRandomTask(t *testing.T, store GlobalRepository) domain.Task {

	user := generateRandomUser(t, store)
	project := generateRandomProject(t, store)

	arg := CreateTaskParams{
		TaskName:  util.RandomName(),
		Budget:    util.RandomBudget(),
		CreatedBy: user.Username,
		TaskOrder: util.RandomInt(0, 100),
		ProjectId: project.Id,
	}

	task, err := store.CreateTask(context.Background(), arg)
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

func GetTaskHelper(t *testing.T, store GlobalRepository) {
	task1 := generateRandomTask(t, store)

	require.NotEmpty(t, task1)

	task2, err := store.GetTask(context.Background(), task1.Id)
	require.NoError(t, err)
	require.NotEmpty(t, task2)
	require.Equal(t, task1.TaskName, task2.TaskName)
	require.Equal(t, task1.Budget, task2.Budget)
	require.Equal(t, task1.CreatedBy, task2.CreatedBy)
	require.Equal(t, task1.TaskOrder, task2.TaskOrder)
	require.Equal(t, task1.ProjectId, task2.ProjectId)
	require.WithinDuration(t, task1.CreatedOn, task2.CreatedOn, time.Second)

}

func GetTaskListHelper(t *testing.T, store GlobalRepository) {
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

func DeleteTaskHelper(t *testing.T, store GlobalRepository) {
	task1 := generateRandomTask(t, store)
	require.NotEmpty(t, task1)

	err := store.DeleteTask(context.Background(), task1.Id)
	require.NoError(t, err)

}

func UpdateTaskHelper(t *testing.T, store GlobalRepository) {
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

func GetUserHelper(t *testing.T, store GlobalRepository) {
	user1 := generateRandomUser(t, store)
	require.NotEmpty(t, user1)

	username := user1.Username

	user2, err := store.GetUser(context.Background(), username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

}
