package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"


	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("Failed to load the config file")
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connet to db: ", err)
	}
	testStore = NewStore(testDB)
	os.Exit(m.Run())

}

func generateRandomUserWithRole(t *testing.T, role string) domain.User {

	hpassword, _ := util.HashPassword(util.RandomString(32))
	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: hpassword,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
		Role:           role,
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user

}



func generateRandomWorksManager(t *testing.T, store GlobalRepository) domain.User {
	return generateRandomUserWithRole(t,util.ROLES[2])
}

func generateRandomClient(t *testing.T, store GlobalRepository) domain.User {
	return generateRandomUserWithRole(t, util.ROLES[1])
}

func generateRandomLocation() domain.Location {
	return domain.Location{
		util.RandomLatitude(), util.RandomLongitude()}
}

func generateRandomProject(t *testing.T) domain.Project {
	resp := generateRandomWorksManager(t, testStore)
	client := generateRandomClient(t, testStore)
	arg := CreateProjectParam{
		ProjectName: util.RandomName(),
		CreatedBy:   resp.Username,
		Client:      client.Username,
		Responsible: resp.Username,
		Address:     util.RandomAddress(),
		Location:    generateRandomLocation(),
	}

	project, err := testStore.CreateProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, project)
	require.Equal(t, arg.ProjectName, project.Projectname)
	require.Equal(t, arg.CreatedBy, project.CreatedBy)
	require.Equal(t, arg.Client, project.Client)
	require.Equal(t, arg.Responsible, project.Responsible)
	require.NotZero(t, project.Id)
	require.NotZero(t, project.CreatedOn)
	return project
}
