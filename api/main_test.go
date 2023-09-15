package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	"github.com/mattn/go-colorable"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var file *os.File

func newTestServerWithEnforcer(t *testing.T, service service.TaskService, enforce bool) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	var err error

	aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

	config, err = util.LoadConfig("..")
	if err != nil {
		logger.Fatal("Can't load the configuration file")
	}
	loaders, err := authorization.Load(config.DBSource, "./config/rbac_model.conf", *logger)
	if err != nil {
		panic(err)
	}
	authorizationService, err := authorization.NewCasbinAuthorization(*loaders)
	policyManagementService, _ := authorization.NewCasbinManagement(*loaders)
	if err != nil {
		panic(err)
	}
	server, err := NewServer(config, service, authorizationService, policyManagementService)
	require.NoError(t, err)
	return server
}

func newTestServer(t *testing.T, service service.TaskService) *Server {
	return newTestServerWithEnforcer(t, service, false)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
	os.RemoveAll(file.Name())
}
