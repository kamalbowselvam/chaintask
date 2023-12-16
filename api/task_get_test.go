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

func TestGetTaskAPI(t *testing.T) {

	project := randomProject(t)
	task := randomTask(project.Client, project.Id, project.CompanyId)

	testCases := []struct {
		name           string
		taskID         int64
		projectID      int64
		companyID      int64
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		setupPolicies  func(t *testing.T, server *Server)
		removePolicies func(t *testing.T, server *Server)
		buildStubs     func(store *mockdb.MockGlobalRepository)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK - Client",
			taskID:    task.Id,
			projectID: task.ProjectId,
			companyID: project.CompanyId,

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
					GetTask(gomock.Any(), task.Id).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				fmt.Println(recorder.Body)
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchTask(t, recorder.Body, task)
			},
		},

		{
			name:      "OK - Responsible",
			taskID:    task.Id,
			projectID: task.ProjectId,
			companyID: project.CompanyId,

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
					GetTask(gomock.Any(), task.Id).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				fmt.Println(recorder.Body)
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchTask(t, recorder.Body, task)
			},
		},

		{
			name:      "NoPolicy",
			taskID:    task.Id,
			projectID: task.ProjectId,
			companyID: project.CompanyId,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, project.Responsible, util.ROLES[2], time.Minute)
			},

			setupPolicies: func(t *testing.T, server *Server) {
				//server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				//server.policies.CreateTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)

			},

			removePolicies: func(t *testing.T, server *Server) {
				//server.policies.RemoveProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
				//server.policies.RemoveTaskPolicies(task.Id, task.ProjectId, task.CreatedBy, task.CompanyId)
			},

			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					GetTask(gomock.Any(), task.Id).
					Times(0).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//fmt.Println(recorder.Body)
				require.Equal(t, http.StatusForbidden, recorder.Code)
				//requiredBodyMatchTask(t, recorder.Body, task)
			},
		},

		// add check for if manager and admin can get the task

		{
			name: "UnauthorizedUser",

			taskID:    task.Id,
			projectID: task.ProjectId,
			companyID: project.CompanyId,

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, "random", "unkown", time.Minute)
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
					GetTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(0).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			taskID:    task.Id,
			projectID: task.ProjectId,
			companyID: project.CompanyId,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
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
					GetTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name:      "NotFound",
			taskID:    task.Id,
			projectID: task.ProjectId,
			companyID: project.CompanyId,

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
					GetTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(1).
					Return(domain.Task{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			taskID:    task.Id,
			projectID: task.ProjectId,
			companyID: project.CompanyId,
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
					GetTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(1).
					Return(domain.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			taskID:    0,
			projectID: task.ProjectId,
			companyID: project.CompanyId,
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
					GetTask(gomock.Any(), gomock.Any()).
					Times(0)
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
			url := fmt.Sprintf("/company/%d/projects/%d/tasks/%d", tc.companyID, tc.projectID, tc.taskID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			tc.setupPolicies(t, server)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
			tc.removePolicies(t, server)
		})

	}

}
