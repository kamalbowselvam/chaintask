package api

import (
	"database/sql"
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
		name           string
		testTask       domain.Task
		testProject    domain.Project
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		setupPolicies  func(t *testing.T, server *Server)
		removePolicies func(t *testing.T, server *Server)
		buildStubs     func(store *mockdb.MockGlobalRepository)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{
			name:        "OK",
			testTask:    task,
			testProject: project,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, project.Client, util.ROLES[1], time.Minute)
			},

			setupPolicies: func(t *testing.T, server *Server) {
				server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.CreateTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			removePolicies: func(t *testing.T, server *Server) {
				server.policies.RemoveProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.RemoveTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)
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
			name:        "NOPolicy",
			testTask:    task,
			testProject: project,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, project.Client, util.ROLES[1], time.Minute)
			},

			setupPolicies: func(t *testing.T, server *Server) {
				//server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				//server.policies.CreateTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			removePolicies: func(t *testing.T, server *Server) {
				// server.policies.RemoveProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				// server.policies.RemoveTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(0).
					Return(nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},

		{
			name:        "KO - Responsible",
			testTask:    task,
			testProject: project,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, project.Responsible, util.ROLES[2], time.Minute)
			},
			setupPolicies: func(t *testing.T, server *Server) {
				server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.CreateTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			removePolicies: func(t *testing.T, server *Server) {
				server.policies.RemoveProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.RemoveTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(0).
					Return(nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},

		{
			name:        "NoAuthorization",
			testTask:    task,
			testProject: project,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, "Random", "unkown", time.Minute)
			},

			setupPolicies: func(t *testing.T, server *Server) {
				server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.CreateTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			removePolicies: func(t *testing.T, server *Server) {
				server.policies.RemoveProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.RemoveTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},

		{
			name:        "NotFound",
			testTask:    task,
			testProject: project,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, project.Client, util.ROLES[1], time.Minute)
			},

			setupPolicies: func(t *testing.T, server *Server) {
				server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.CreateTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			removePolicies: func(t *testing.T, server *Server) {
				server.policies.RemoveProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.RemoveTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		{
			name:        "Internal Error",
			testTask:    task,
			testProject: project,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, project.Client, util.ROLES[1], time.Minute)
			},

			setupPolicies: func(t *testing.T, server *Server) {
				server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.CreateTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			removePolicies: func(t *testing.T, server *Server) {
				server.policies.RemoveProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.RemoveTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "InvalidID",
			testTask: domain.Task{
				Id:        0,
				CreatedBy: project.Client,
				CompanyId: task.CompanyId,
				ProjectId: task.ProjectId,
			},
			testProject: project,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, project.Client, util.ROLES[1], time.Minute)
			},

			setupPolicies: func(t *testing.T, server *Server) {
				server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.CreateTaskPolicies(0, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			removePolicies: func(t *testing.T, server *Server) {
				server.policies.RemoveProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				server.policies.RemoveTaskPolicies(0, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/company/%d/projects/%d/tasks/%d", tc.testTask.CompanyId, tc.testTask.ProjectId, tc.testTask.Id)

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			tc.setupPolicies(t, server)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
			tc.removePolicies(t, server)

		})

	}

}
