package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type updateTaskParams struct {
	TaskName  string `json:"task_name" binding:"required"`
	UpdatedBy string `swaggerignore:"true"`
	TaskOrder int64  `json:"task_order" binding:"required,number"`
	ProjectId int64  `json:"project_id" binding:"required,number"`
}

type eqUpdateTaskParamsMatcher struct {
	arg db.UpdateTaskParams
}

func (e eqUpdateTaskParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateTaskParams)
	if !ok {
		return false
	}

	if !(e.arg.Budget.Equal(arg.Budget)) {
		return false
	}

	if e.arg.UpdatedBy != arg.UpdatedBy {
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

	eparam := updateTaskParams{TaskName: e.arg.TaskName,
		UpdatedBy: e.arg.UpdatedBy,
		ProjectId: *e.arg.ProjectId,
		TaskOrder: e.arg.TaskOrder}

	argparam := updateTaskParams{TaskName: arg.TaskName,
		UpdatedBy: arg.UpdatedBy,
		ProjectId: *arg.ProjectId,
		TaskOrder: arg.TaskOrder}

	return reflect.DeepEqual(eparam, argparam)
}

func (e eqUpdateTaskParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqUpdateTaskParams(arg db.UpdateTaskParams) gomock.Matcher {
	return eqUpdateTaskParamsMatcher{arg}
}

func TestUpdateTaskAPI(t *testing.T) {

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
				"done":       task.Done,
				"project_id": task.ProjectId,
				"task_order": task.TaskOrder,
				"version":    task.Version,
				"rating":     task.Rating,
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
				
					arg := db.UpdateTaskParams{
						TaskName:  task.TaskName,
						UpdatedBy: task.CreatedBy,
						Budget:    task.Budget,
						ProjectId: &task.ProjectId,
						TaskOrder: task.TaskOrder,
						Done:      task.Done,
						Version:   &task.Version,
						Rating:    &task.Rating,
					}
				
				store.EXPECT().
					UpdateTask(gomock.Any(), EqUpdateTaskParams(arg)).
					Times(1).
					Return(task, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//				requiredBodyMatchTask(t, recorder.Body, task)
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

			server := newTestServerWithEnforcer(t, store, true)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/company/%d/projects/%d/tasks/%d", tc.testTask.CompanyId, tc.testTask.ProjectId, tc.testTask.Id)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			tc.setupPolicies(t, server)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
			tc.removePolicies(t, server)

		})

	}

}
