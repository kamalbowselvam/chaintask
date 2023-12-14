package api

import (
	"bytes"
	"database/sql"
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

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	mockdb "github.com/kamalbowselvam/chaintask/mock"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

type createTaskParams struct {
	TaskName  string `json:"task_name" binding:"required"`
	CreatedBy string `swaggerignore:"true"`
	TaskOrder int64  `json:"task_order" binding:"required,number"`
	ProjectId int64  `json:"project_id" binding:"required,number"`
}

type eqCreateTaskParamsMatcher struct {
	arg db.CreateTaskParams
}

func (e eqCreateTaskParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateTaskParams)
	if !ok {
		return false
	}

	if !(e.arg.Budget.Equal(arg.Budget)) {
		return false
	}

	if e.arg.CreatedBy != arg.CreatedBy {
		return false
	}

	if e.arg.TaskName != arg.TaskName {
		return false
	}

	// this conversion without budget is done to avoid doing deep comparison of Decimal (Budget object)
	// which causes error due to Exponentional

	// If given argument is Decimal(473,1) --> Which is 4730
	// Request param is unmarshalled as Decimal(4730,0) --> Which is also 4730
	// but deepequal throws error for the object

	eparam := createTaskParams{TaskName: e.arg.TaskName,
		CreatedBy: e.arg.CreatedBy,
		ProjectId: *e.arg.ProjectId,
		TaskOrder: e.arg.TaskOrder}

	argparam := createTaskParams{TaskName: arg.TaskName,
		CreatedBy: arg.CreatedBy,
		ProjectId: *arg.ProjectId,
		TaskOrder: arg.TaskOrder}

	return reflect.DeepEqual(eparam, argparam)
}

func (e eqCreateTaskParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqCreateTaskParams(arg db.CreateTaskParams) gomock.Matcher {
	return eqCreateTaskParamsMatcher{arg}
}

func TestCreateTaskAPI(t *testing.T) {

	project := randomProject(t)
	task := randomTask(project.Client, project.Id, project.CompanyId)

	testCases := []struct {
		name           string
		body           gin.H
		testTask       domain.Task
		testProject    domain.Project
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		setupPolicies  func(t *testing.T, server *Server)
		removePolicies func(t *testing.T, server *Server)
		buildStubs     func(store *mockdb.MockGlobalRepository)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK - Client",
			testTask:    task,
			testProject: project,
			body: gin.H{
				"task_name":  task.TaskName,
				"budget":     task.Budget,
				"project_id": task.ProjectId,
				"task_order": task.TaskOrder,
			},

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
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: project.Client,
					ProjectId: &task.ProjectId,
					TaskOrder: task.TaskOrder,
				}

				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(task, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requiredBodyMatchTask(t, recorder.Body, task)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:        "OK - Reponsible",
			testTask:    task,
			testProject: project,

			body: gin.H{
				"task_name":  task.TaskName,
				"budget":     task.Budget,
				"project_id": task.ProjectId,
				"task_order": task.TaskOrder,
			},

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
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: project.Responsible,
					ProjectId: &task.ProjectId,
					TaskOrder: task.TaskOrder,
				}

				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(task, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requiredBodyMatchTask(t, recorder.Body, task)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:        "UnauthorizedUser",
			testTask:    task,
			testProject: project,

			body: gin.H{
				"task_name":  task.TaskName,
				"budget":     task.Budget,
				"project_id": task.ProjectId,
				"task_order": task.TaskOrder,
			},

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, "Random", "unknown", time.Minute)
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
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: task.CreatedBy,
					ProjectId: &task.ProjectId,
					TaskOrder: task.TaskOrder,
				}

				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(0).
					Return(task, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requiredBodyMatchTask(t, recorder.Body, task)
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},

		{
			name:        "NoAuthorization",
			testTask:    task,
			testProject: project,

			body: gin.H{
				"task_name":  task.TaskName,
				"budget":     task.Budget,
				"project_id": task.ProjectId,
				"task_order": task.TaskOrder,
			},

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
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: task.CreatedBy,
					ProjectId: &task.ProjectId,
					TaskOrder: task.TaskOrder,
				}

				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(0).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name:        "InternalError",
			testTask:    task,
			testProject: project,

			body: gin.H{
				"task_name":  task.TaskName,
				"budget":     task.Budget,
				"project_id": task.ProjectId,
				"task_order": task.TaskOrder,
			},

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
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: task.CreatedBy,
					ProjectId: &task.ProjectId,
					TaskOrder: task.TaskOrder,
				}

				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(domain.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name:        "InvalidData",
			testTask:    task,
			testProject: project,

			body: gin.H{
				"task":       task.TaskName,
				"budget":     task.Budget,
				"project_id": task.ProjectId,
				"task_order": task.TaskOrder,
			},

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
				/*
					arg := db.CreateTaskParams{
						TaskName:  task.TaskName,
						Budget:    task.Budget,
						CreatedBy: task.CreatedBy,
						ProjectId: &task.ProjectId,
						TaskOrder: task.TaskOrder,
					}
				*/

				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(0).
					Return(domain.Task{}, sql.ErrConnDone)
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/company/%d/projects/%d/tasks/", tc.testTask.CompanyId, tc.testTask.ProjectId)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			tc.setupPolicies(t, server)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
			tc.removePolicies(t, server)

		})

	}

}

func randomProject(t *testing.T) domain.Project {

	client, _ := randomUser(t, util.ROLES[0])
	responsible, _ := randomUser(t, util.ROLES[1])

	return domain.Project{
		Id:                   util.RandomInt(1, 1000),
		Projectname:          util.RandomName(),
		CreatedOn:            time.Now(),
		CreatedBy:            util.DEFAULT_SUPER_ADMIN,
		Longitude:            util.RandomLongitude(),
		Latitude:             util.RandomLatitude(),
		Address:              util.RandomAddress(),
		Client:               client.Username,
		Responsible:          responsible.Username,
		Budget:               float64(util.RandomInt(1, 100)),
		CompletionPercentage: float64(util.RandomInt(1, 100)),
		CompanyId:            util.RandomInt(1, 100),
	}
}

//func randomProjectWithinCompany(client string, responsible string, companyId int64) domain.Project {
//	project := randomProject(client, responsible)
//	project.CompanyId = companyId
//	return project
//}

func randomTask(username string, projectId int64, companyId int64) domain.Task {

	name := util.RandomName()
	budget := util.RandomBudget()

	return domain.Task{
		Id:        util.RandomInt(1, 1000),
		TaskName:  name,
		Budget:    budget,
		CreatedBy: username,
		UpdatedBy: username,
		ProjectId: projectId,
		TaskOrder: util.RandomInt(1, 10),
		CompanyId: companyId,
		Rating:    util.RandomInt(1, 5),
		Version:   util.RandomInt(1, 10),
		Done:      true,
	}
}

func requiredBodyMatchTask(t *testing.T, body *bytes.Buffer, task domain.Task) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTask domain.Task

	err = json.Unmarshal(data, &gotTask)
	require.NoError(t, err)
	require.Equal(t, task.CreatedOn.UnixMilli(), gotTask.CreatedOn.UnixMilli())

}
