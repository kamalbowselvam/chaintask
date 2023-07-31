package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	mockdb "github.com/kamalbowselvam/chaintask/mock"
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

	if e.arg.TaskName  != arg.TaskName {
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
	task := randomTask()

	testCases := []struct {
		name         string
		taskID       int64
		buidStubs    func(store *mockdb.MockGlobalRepository)
		checkReponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			taskID: task.Id,
			buidStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					GetTask(gomock.Any(), task.Id).
					Times(1).
					Return(task, nil)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchTask(t, recorder.Body, task)
			},
		},

		{
			name:   "NotFound",
			taskID: task.Id,
			buidStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					GetTask(gomock.Any(), task.Id).
					Times(1).
					Return(domain.Task{}, errors.New("id not found"))
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockGlobalRepository(ctrl)
			tc.buidStubs(store)

			taskService := NewTaskService(store)
			taskHandler := rest.NewHttpHandler(taskService)
			router := gin.New()
			router.GET("/tasks/:id", taskHandler.GetTask)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/tasks/%d", task.Id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)
			tc.checkReponse(t, recorder)

		})

	}

}

func TestCreateTaskAPI(t *testing.T) {
	task := randomTask()
	t.Log(task)
	testCases := []struct {
		name          string
		body 		  gin.H
		gtask         domain.Task
		buildStubs    func(store *mockdb.MockGlobalRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			gtask: task,
			body: gin.H{
				"taskname": task.TaskName,
				"createdBy": task.CreatedBy,
				"budget": task.Budget,
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				arg := db.CreateTaskParams{
					TaskName: task.TaskName,
					Budget: task.Budget,
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

			taskService := NewTaskService(store)
			taskHandler := rest.NewHttpHandler(taskService)

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

func randomTask() domain.Task {

	name := util.RandomName()
	budget := util.RandomBudget()
	createdBy := util.RandomName()
	task := domain.NewTask(name, budget, createdBy)
	//task.Id = util.RandomInt(1, 100)
	return task
}

func requiredBodyMatchTask(t *testing.T, body *bytes.Buffer, task domain.Task) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTask domain.Task

	err = json.Unmarshal(data, &gotTask)
	require.NoError(t, err)
	require.Equal(t, task.CreatedOn.UnixMilli(), gotTask.CreatedOn.UnixMilli())

}
