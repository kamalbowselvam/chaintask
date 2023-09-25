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

	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	mockdb "github.com/kamalbowselvam/chaintask/mock"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

type eqCreateTaskParamsMatcher struct {
	arg db.CreateTaskParams
}

func (e eqCreateTaskParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateTaskParams)
	if !ok {
		return false
	}

	if e.arg.Budget != arg.Budget {
		return false
	}

	if e.arg.CreatedBy != arg.CreatedBy {
		return false
	}

	if e.arg.TaskName != arg.TaskName {
		return false
	}

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateTaskParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqCreateTaskParams(arg db.CreateTaskParams) gomock.Matcher {
	return eqCreateTaskParamsMatcher{arg}
}

func TestGetTaskAPI(t *testing.T) {
	client, _ := randomUser(t, util.ROLES[0])
	responsible, _ := randomUser(t, util.ROLES[1])
	project := randomProject(client.Username, responsible.Username)
	task := randomTask(client.Username, project.Id)

	testCases := []struct {
		name          string
		taskID        int64
		projectID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockGlobalRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			taskID:    task.Id,
			projectID: task.ProjectId,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, client.Username, client.Role, time.Minute)
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
			name:      "UnauthorizedUser",
			taskID:    task.Id,
			projectID: task.ProjectId,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					GetTask(gomock.Any(), gomock.Eq(task.Id)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			taskID:    task.Id,
			projectID: task.ProjectId,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, client.Username, client.Role, time.Minute)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, client.Username, client.Role, time.Minute)
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
			projectID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, client.Username, client.Role, time.Minute)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/projects/%d/tasks/%d", tc.projectID, tc.taskID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}

}

func TestCreateTaskAPI(t *testing.T) {
	user, _ := randomUser(t, util.ROLES[1])
	responsible, _ := randomUser(t, util.ROLES[2])
	project := randomProject(user.Username, responsible.Username)
	task := randomTask(user.Username, project.Id)
	config := loadConfig()
	authorizationLoaders := generateLoader(*config)

	testCases := []struct {
		name          string
		body          gin.H
		gtask         domain.Task
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		setupAuthorization func(t *testing.T, authorizationLoaders *authorization.Loaders)
		buildStubs    func(store *mockdb.MockGlobalRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			gtask: task,
			body: gin.H{
				"taskname":  task.TaskName,
				"createdBy": task.CreatedBy,
				"budget":    task.Budget,
				"projectId": task.ProjectId,
				"taskOrder": task.TaskOrder,
			},

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},
			setupAuthorization: func(t *testing.T, authorizationLoaders *authorization.Loaders){
				AddAuthorization(t, *authorizationLoaders, user.Username, fmt.Sprintf("/projects/%d/tasks/", task.ProjectId), http.MethodPost)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: task.CreatedBy,
					ProjectId: task.ProjectId,
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

			url := fmt.Sprintf("/projects/%d/tasks/", tc.gtask.ProjectId)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			tc.setupAuthorization(t, authorizationLoaders)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})

	}

}

func randomProject(client string, responsible string) domain.Project {
	return domain.Project{
		Id:          util.RandomInt(1, 1000),
		Projectname: util.RandomName(),
		CreatedOn:   time.Now(),
		CreatedBy:   util.DEFAULT_ADMIN,
		Location: domain.Location{
			util.RandomLatitude(),
			util.RandomLongitude(),
		},
		Address:              util.RandomAddress(),
		Client:               client,
		Responsible:          responsible,
		Budget:               float64(util.RandomInt(1, 100)),
		CompletionPercentage: float64(util.RandomInt(1, 100)),
	}
}

func randomTask(username string, projectId int64) domain.Task {

	name := util.RandomName()
	budget := util.RandomBudget()

	return domain.Task{
		Id:        util.RandomInt(1, 1000),
		TaskName:  name,
		Budget:    budget,
		CreatedBy: username,
		ProjectId: projectId,
		TaskOrder: util.RandomInt(1, 10),
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
