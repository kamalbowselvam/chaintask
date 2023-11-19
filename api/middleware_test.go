package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	mockdb "github.com/kamalbowselvam/chaintask/mock"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func addAuthentification(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authentificationType string,
	username string,
	user_role string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(username, user_role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authentificationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func AddAuthorization(
	t *testing.T,
	authorizationLoaders authorization.Loaders,
	username string,
	resource string,
	rights string,
){
	authorizationLoaders.Enforcer.EnableEnforce(true)
	// Beware, wipes all entries from casbin DB
	// this line should actually be _, err := authorizationLoaders.Enforcer.RemoveFilteredNamedPolicy("p", 0, "");
	// but see with Kamal if it affects his workflow
	_, err := authorizationLoaders.Enforcer.RemovePolicies([][]string{{"*"}})
	if err!= nil{
		logger.Warn(err.Error())
	}
	rules := [][]string{
		{username, resource, rights},
	}
	_, err  = authorizationLoaders.Enforcer.AddPoliciesEx(rules)
	require.NoError(t, err)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, "user", "USER", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, "unsupported", "user", "USER", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, "", "user", "USER", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, "user", "USER", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockGlobalRepository(ctrl)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestAuthorizationMiddleware(t *testing.T) {
	admin, _ := randomUser(t, util.ROLES[3])
	client, _ := randomUser(t, util.ROLES[1])
	user, _ := randomUser(t, util.ROLES[1])
	responsible, _ := randomUser(t, util.ROLES[2])
	project := randomProject(client.Username, responsible.Username)
	companyId := project.CompanyId
	task := randomTask(client.Username, project.Id)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		gtask         domain.Task
		buildStubs    func(store *mockdb.MockGlobalRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "NOK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.UserRole, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					ProjectId: task.ProjectId,
					TaskOrder: task.TaskOrder,
				}

				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(0).
					Return(task, nil)

			},
			body: gin.H{
				"task_name":  task.TaskName,
				"budget":    task.Budget,
				"task_order": task.TaskOrder,
				"project_id": task.ProjectId,
			},
			gtask: task,
		},
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, admin.Username, admin.UserRole, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: admin.Username,
					ProjectId: task.ProjectId,
					TaskOrder: task.TaskOrder,
				}

				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(task, nil)

			},
			body: gin.H{
				"task_name":  task.TaskName,
				"budget":    task.Budget,
				"task_order": task.TaskOrder,
				"project_id": task.ProjectId,
			},
			gtask: task,
		},
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, client.Username, client.UserRole, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: client.Username,
					ProjectId: task.ProjectId,
					TaskOrder: task.TaskOrder,
				}

				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(task, nil)

			},
			body: gin.H{
				"task_name":  task.TaskName,
				"budget":    task.Budget,
				"task_order": task.TaskOrder,
				"project_id": task.ProjectId,
			},
			gtask: task,
		},
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, responsible.Username, responsible.UserRole, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: responsible.Username,
					ProjectId: task.ProjectId,
					TaskOrder: task.TaskOrder,
				}

				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(task, nil)

			},
			body: gin.H{
				"task_name":  task.TaskName,
				"budget":    task.Budget,
				"task_order": task.TaskOrder,
				"project_id": task.ProjectId,
			},
			gtask: task,
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
			server.policies.CreateAdminPolicies(admin.Username, admin.CompanyId)
			server.policies.CreateProjectPolicies(project.Id, project.Client, project.Responsible, project.CompanyId)
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)

			require.NoError(t, err)

			url := fmt.Sprintf("/company/%d/projects/%d/tasks/", companyId, tc.gtask.ProjectId)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
