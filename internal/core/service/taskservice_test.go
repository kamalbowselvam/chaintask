package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kamalbowselvam/chaintask/internal/core/domain"
	"github.com/kamalbowselvam/chaintask/internal/handlers/rest"
	mockdb "github.com/kamalbowselvam/chaintask/internal/mock"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func TestGetTaskAPI(t *testing.T) {
	task := randomTask()

	testCases := []struct {
		name         string
		taskID       int64
		buidStubs    func(store *mockdb.MockTaskRepository)
		checkReponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			taskID: task.Id,
			buidStubs: func(store *mockdb.MockTaskRepository) {
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
			buidStubs: func(store *mockdb.MockTaskRepository) {
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
			store := mockdb.NewMockTaskRepository(ctrl)
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

	testCases := []struct {
		name          string
		gtask         domain.Task
		buildStubs    func(store *mockdb.MockTaskRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			gtask: task,
			buildStubs: func(store *mockdb.MockTaskRepository) {
				store.EXPECT().
					SaveTask(gomock.Any(), gomock.Any()).
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
			store := mockdb.NewMockTaskRepository(ctrl)
			tc.buildStubs(store)

			taskService := NewTaskService(store)
			taskHandler := rest.NewHttpHandler(taskService)

			router := gin.New()
			router.POST("/tasks/", taskHandler.CreateTask)
			recorder := httptest.NewRecorder()

			body, err := json.Marshal(tc.gtask)
			//name := tc.gtask.Name
			//budget := tc.gtask.Budget
			//createdby := tc.gtask.CreatedBy

			//b := fmt.Sprintf("{\"name\": \"%s\", \"budget\": \"%f\", \"createdBy\": \"%s\"}",name,budget,createdby )
			//body := []byte(b)

			require.NoError(t, err)
			//fmt.Println(body)
			jsonbody := bytes.NewReader(body)

			t.Log(string(body))
			url := "/tasks/"

			request, err := http.NewRequest(http.MethodPost, url, jsonbody)
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