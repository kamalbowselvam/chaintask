package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kamalbowselvam/chaintask/domain"
	mockdb "github.com/kamalbowselvam/chaintask/mock"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/stretchr/testify/require"
)

func randomProject(t *testing.T, client domain.User, responsible domain.User, admin domain.User) (project domain.Project) {
	
	project = domain.Project{
		Id : util.RandomInt(0,1000),
		Projectname: util.RandomString(10),
		CreatedBy:   admin.Username,
		Location:    domain.Location{util.RandomLatitude(), util.RandomLongitude()},
		Responsible: responsible.Username,
		Client:      client.Username,
		Address:     util.RandomAddress(),
	}
	return
}

func TestCreateProjectAPI(t *testing.T) {
	client, _ := randomUser(t, util.ROLES[1])
	responsible, _ := randomUser(t, util.ROLES[2])
	admin, _ := randomUser(t, util.ROLES[3])

	project := randomProject(t, client, responsible, admin)

	testCases := []struct {
		name          string
		projectID     int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockGlobalRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			projectID: project.Id,
			body: gin.H{
				"projectname":  project.Projectname,
				"createdBy": project.CreatedBy,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, admin.Username, admin.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockGlobalRepository) {
				store.EXPECT().
					CreateProject(gomock.Any(), gomock.Any()).
					Times(1).
					Return(project, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchProject(t, recorder.Body, project)
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
			authRoutes.POST("/projects/", taskHandler.CreateProject)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/projects/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			tc.setupAuth(t, request, taskHandler.tokenMaker)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})

	}



}
