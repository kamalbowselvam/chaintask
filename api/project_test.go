package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kamalbowselvam/chaintask/domain"
	mockdb "github.com/kamalbowselvam/chaintask/mock"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)


func requiredBodyMatchProject(t *testing.T, body *bytes.Buffer, task domain.Project) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTask domain.Project

	err = json.Unmarshal(data, &gotTask)
	require.NoError(t, err)
	require.Equal(t, task.CreatedOn.UnixMilli(), gotTask.CreatedOn.UnixMilli())

}



func TestCreateProjectAPI(t *testing.T) {
	client, _ := randomUser(t, util.ROLES[1])
	responsible, _ := randomUser(t, util.ROLES[2])
	admin, _ := randomUser(t, util.ROLES[3])

	project := randomProject(client.Username, responsible.Username)

	testCases := []struct {
		name          string
		projectID     int64
		companyID     int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockGlobalRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			projectID: project.Id,
			companyID: project.CompanyId,
			body: gin.H{
				"projectname":  project.Projectname,
				"createdBy": project.CreatedBy,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, admin.Username, admin.UserRole, time.Minute)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					CreateProject(gomock.Any(), gomock.Any()).
					Times(1).
					Return(project, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchProject(t, recorder.Body, project)
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

			server := newTestServerWithEnforcer(t, store, true)
			server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible,  project.CompanyId)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/company/%d/projects/", tc.companyID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

