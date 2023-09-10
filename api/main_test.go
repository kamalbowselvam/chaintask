package api

import (
	"os"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, service service.TaskService) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	adapter := fileadapter.NewAdapter("../tests/fake_policy.csv")
	enforcer, err := casbin.NewEnforcer("../config/rbac_model.conf", adapter)
	loaders := authorization.FakeLoader{
		Adapter:  adapter,
		Enforcer: enforcer,
	}
	if err != nil {
		panic(err)
	}
	authorizationService, err := authorization.NewFakeCasbinAuthorization(loaders)
	require.NoError(t, err)
	policyManagementService, err := authorization.NewFakeCasbinManagement(loaders)
	require.NoError(t, err)
	server, err := NewServer(config, service, authorizationService, policyManagementService)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
