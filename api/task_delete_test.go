package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kamalbowselvam/chaintask/domain"
	mockdb "github.com/kamalbowselvam/chaintask/mock"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func TestDeleteTaskAPI(t *testing.T) {


	project := randomProject(t)
	task := randomTask(project.Client, project.Id, project.CompanyId)

	testCases := []struct {
		name               string
		gtask              domain.Task
		setupAuth          func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs         func(store *mockdb.MockGlobalRepository)
		checkResponse      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{
			name:  "OK",
			gtask: task,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, project.Client, util.ROLES[1], time.Minute)
			},

			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(1).
					Return(nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
			},
		},

		{
			name:  "KO - Responsible",
			gtask: task,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, project.Responsible, util.ROLES[2], time.Minute)
			},

			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(0).
					Return(nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
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

			server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)

			url := fmt.Sprintf("/company/%d/projects/%d/tasks/%d", task.CompanyId, tc.gtask.ProjectId, task.Id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})

	}

}
