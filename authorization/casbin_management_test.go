package authorization

import (
	"net/http"
	"testing"

	"github.com/kamalbowselvam/chaintask/util"
	"github.com/mattn/go-colorable"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCreatePolicies(t *testing.T) {

	aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	config, err := util.LoadConfig("..")
	if err != nil {
		logger.Fatal("Can't load the configuration file")
	}
	loaders, err := Load(config.DBSource, "./config/rbac_model.conf", *logger)
	if err != nil {
		panic(err)
	}
	policyManagementService, err := NewCasbinManagement(*loaders)
	if err != nil {
		panic(err)
	}
	name := util.RandomName()
	err = policyManagementService.CreateAdminPolicies(name)
	require.NoError(t, err)
	right, err2 := loaders.Enforcer.Enforce(name, "/projects", http.MethodPost)
	require.True(t, right)
	require.NoError(t, err2)
}
