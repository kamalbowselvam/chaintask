package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	"github.com/mattn/go-colorable"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	if err != nil {
		logger.Fatal("Failed to load the config file")
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Fatal("cannot connet to db: ", zap.Error(err))
	}



	testStore = NewStore(testDB)
	os.Exit(m.Run())

}

func generateRandomCompany(t *testing.T) domain.Company {
	arg := CreateCompanyParams{
		CompanyName: util.RandomName(),
		Address: util.RandomAddress(),
		CreatedBy: util.RandomName(),
	}

	company, err := testStore.CreateCompany(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, company)
	require.Equal(t, arg.CompanyName, company.CompanyName)
	require.Equal(t, arg.Address, company.Address)
	require.Equal(t, arg.CreatedBy, company.CreatedBy)
	require.Equal(t, arg.CreatedBy, company.UpdatedBy)

	require.NotZero(t, company.Id)
	require.NotZero(t, company.CreatedOn)
	require.NotZero(t, company.UpdatedOn)

	return company
}

func generateRandomUserWithRoleAndCompany(t *testing.T, role string, company int64) domain.User{
	user := generateRandomUserWithRole(t, role);
	user.CompanyId = company
	return user;

}
func generateRandomUserWithRole(t *testing.T, role string) domain.User {

	hpassword, _ := util.HashPassword(util.RandomString(32))
	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: hpassword,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
		UserRole:       role,
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user

}

func generateRandomWorksManager(t *testing.T) domain.User {
	return generateRandomUserWithRole(t, "RESPONSIBLE")
}

func generateRandomWorksManagerWithinCompany(t *testing.T, company int64) domain.User {
	return generateRandomUserWithRoleAndCompany(t, "RESPONSIBLE", company)
}

/*
func generateRandomClient(t *testing.T) domain.User {
	return generateRandomUserWithRole(t, "CLIENT")
}*/

func generateRandomClientWithinCompany(t *testing.T, company int64) domain.User {
	return generateRandomUserWithRoleAndCompany(t, "CLIENT", company)
}

func generateRandomLocation() domain.Location {
	return domain.Location{
		util.RandomLatitude(), util.RandomLongitude()}
}

func generateRandomProject(t *testing.T) domain.Project {
	company := generateRandomCompany(t)
	resp := generateRandomWorksManagerWithinCompany(t, company.Id)
	client := generateRandomClientWithinCompany(t, company.Id)
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
	require.NotZero(t, project.CompanyId)
	return project
}
