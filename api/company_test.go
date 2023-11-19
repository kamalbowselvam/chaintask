package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	mockdb "github.com/kamalbowselvam/chaintask/mock"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func randomCompany(t *testing.T) domain.Company{
	return domain.Company{
		CompanyName: util.RandomString(10),
		Address: util.RandomAddress(),
	}
}

func generateRandomCompany(t *testing.T) domain.Company {
	arg := db.CreateCompanyParams{
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

type eqCreateCompanyParamsMatcher struct {
	arg db.CreateCompanyParams
}

func (e eqCreateCompanyParamsMatcher) Matches(x interface{}) bool {
	return reflect.DeepEqual(e.arg, x)
}

func (e eqCreateCompanyParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqCreateCompanyParams(arg db.CreateCompanyParams) gomock.Matcher {
	return eqCreateCompanyParamsMatcher{arg}
}

func requiredBodyMatchCompany(t *testing.T, body *bytes.Buffer, Company domain.Company) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCompany domain.Company

	err = json.Unmarshal(data, &gotCompany)
	require.NoError(t, err)
	require.Equal(t, Company.CreatedOn.UnixMilli(), gotCompany.CreatedOn.UnixMilli())

}


func TestCreateCompanyAPI(t *testing.T) {
	admin, _ := randomUser(t, util.ROLES[3])
	config := loadConfig()
	company := randomCompany(t)
	company.CreatedBy = admin.Username
	authorizationLoaders := generateLoader(*config)

	testCases := []struct {
		name          string
		body          gin.H
		gtask         domain.Company
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		setupAuthorization func(t *testing.T, authorizationLoaders *authorization.Loaders)
		buildStubs    func(store *mockdb.MockGlobalRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			gtask: company,
			body: gin.H{
				"companyname":  company.CompanyName,
				"address" : company.Address,
			},

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, admin.Username, admin.UserRole, time.Minute)
			},
			setupAuthorization: func(t *testing.T, authorizationLoaders *authorization.Loaders){
				AddAuthorization(t, *authorizationLoaders, admin.Username, "/companies*", http.MethodPost)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				arg := db.CreateCompanyParams{
					CompanyName:  company.CompanyName,
					Address:    company.Address,
					CreatedBy: company.CreatedBy,
				}

				store.EXPECT().
					CreateCompany(gomock.Any(), EqCreateCompanyParams(arg)).
					Times(1).
					Return(company, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requiredBodyMatchCompany(t, recorder.Body, company)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()
			store := mockdb.NewMockGlobalRepository(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)

			require.NoError(t, err)

			url := "/companies/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			tc.setupAuthorization(t, authorizationLoaders)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})

	}

}