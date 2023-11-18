package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func loadConfig() *util.Config {
	config, err := util.LoadConfig("..")
	if err != nil {
		panic(err)
	}
	return &config
}

func generateLoader(config util.Config) *authorization.Loaders{
	//aa := zap.NewDevelopmentEncoderConfig()
	//aa.EncodeLevel = zapcore.CapitalColorLevelEncoder

	//logger := zap.New(zapcore.NewCore(
	//	zapcore.NewConsoleEncoder(aa),
	//	zapcore.AddSync(colorable.NewColorableStdout()),
	//	zapcore.DebugLevel,
	//))
	loaders, err := authorization.Load(config.DBSource, "./config/rbac_model.conf")
	if err != nil {
		panic(err)
	}
	return loaders
}

func newTestServerWithEnforcerAndLoaders(t *testing.T, service service.TaskService, enforce bool, loaders authorization.Loaders, config util.Config) *Server{
	loaders.Enforcer.EnableEnforce(enforce)
	authorizationService, err := authorization.NewCasbinAuthorization(loaders)
	policyManagementService, _ := authorization.NewCasbinManagement(loaders)
	if err != nil {
		panic(err)
	}
	server, err := NewServer(config, service, authorizationService, policyManagementService)
	require.NoError(t, err)
	return server
}

func newTestServerWithEnforcer(t *testing.T, service service.TaskService, enforce bool) *Server {
	config := loadConfig()
	loaders := generateLoader(*config)
	return newTestServerWithEnforcerAndLoaders(t, service, enforce, *loaders, *config)
}

func newTestServer(t *testing.T, service service.TaskService) *Server {
	return newTestServerWithEnforcer(t, service, false)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
