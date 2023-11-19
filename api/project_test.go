package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
	"github.com/mattn/go-colorable"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


func requiredBodyMatchProject(t *testing.T, body *bytes.Buffer, task domain.Project) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTask domain.Project

	err = json.Unmarshal(data, &gotTask)
	require.NoError(t, err)
	require.Equal(t, task.CreatedOn.UnixMilli(), gotTask.CreatedOn.UnixMilli())

}



func TestCreateProjectAPI(t *testing.T) {
	config, err := util.LoadConfig("../")
	aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	if err != nil {
		logger.Fatal("Failed to load the config file")
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Fatal("cannot connet to db: ", zap.Error(err))
	}

	testStore = db.NewStore(testDB)
	company := generateRandomCompany(t)
	client, _ := randomUserWithinCompany(t, util.ROLES[1], company.Id)
	responsible, _ := randomUserWithinCompany(t, util.ROLES[2], company.Id)
	admin, _ := randomUserWithinCompany(t, util.ROLES[3], company.Id)

	project := randomProjectWithinCompany(client.Username, responsible.Username, company.Id)

	testCases := []struct {
		name          string
		projectID     int64
		companyID     int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		setupAuthorization func(t *testing.T, authorizationLoaders *authorization.Loaders)
		buildStubs    func(store *mockdb.MockGlobalRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			projectID: project.Id,
			companyID: project.CompanyId,
			body: gin.H{
				"projectname":  project.Projectname,
				"createdBy": project.CreatedBy,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentification(t, request, tokenMaker, authorizationTypeBearer, admin.Username, admin.UserRole, time.Minute)
			},
			setupAuthorization: func(t *testing.T, authorizationLoaders *authorization.Loaders){
				AddAuthorization(t, *authorizationLoaders, admin.Username, fmt.Sprintf("/company/%d/projects", admin.CompanyId), http.MethodPost)
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

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)

			require.NoError(t, err)
			url := fmt.Sprintf("/company/%d/projects/", tc.companyID)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

