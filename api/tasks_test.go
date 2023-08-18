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
	user, _ := randomUser(t)
	task := randomTask(user.Username)

	testCases := []struct {
		name          string
		taskID        int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockGlobalRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			taskID: task.Id,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role,  time.Minute)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					GetTask(gomock.Any(), task.Id).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name:   "UnauthorizedUser",
			taskID: task.Id,
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
			name:   "NoAuthorization",
			taskID: task.Id,
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
			name:   "NotFound",
			taskID: task.Id,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
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
			name:   "InternalError",
			taskID: task.Id,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
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
			name:   "InvalidID",
			taskID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
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

			taskHandler := NewTestHandler(t, store)
			router := gin.New()
			authRoutes := router.Group("/").Use(AuthMiddleware(taskHandler.tokenMaker))
			authRoutes.GET("/tasks/:id", taskHandler.GetTask)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/tasks/%d", tc.taskID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(t, request, taskHandler.tokenMaker)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})

	}

}

func TestCreateTaskAPI(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(user.Username)

	t.Log(task)
	testCases := []struct {
		name          string
		body          gin.H
		gtask         domain.Task
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
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				arg := db.CreateTaskParams{
					TaskName:  task.TaskName,
					Budget:    task.Budget,
					CreatedBy: task.CreatedBy,
				}
				store.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(task, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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

			taskHandler := NewTestHandler(t, store)

			router := gin.New()
			router.POST("/tasks/", taskHandler.CreateTask)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/tasks/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)

		})

	}

}

func randomTask(username string) domain.Task {

	name := util.RandomName()
	budget := util.RandomBudget()

	//task.Id = util.RandomInt(1, 100)
	return domain.Task{
		Id:        util.RandomInt(1, 1000),
		TaskName:  name,
		Budget:    budget,
		CreatedBy: username,
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
