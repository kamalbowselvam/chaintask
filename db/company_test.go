package db

import (
	"context"
	"testing"

	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func generateRandomCompany(t *testing.T) domain.Company {
	arg := CreateCompanyParams{
		CompanyName: util.RandomName(),
		Address: util.RandomString(100),
		CreatedBy: util.RandomString(10),
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

func TestCreateCompany(t *testing.T) {
	generateRandomCompany(t)
}