package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func NewTestHandler(t *testing.T, store db.GlobalRepository) *HttpHandler {

	service := service.NewTaskService(store)

	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	require.NoError(t,err)
	handler := NewHttpHandler(service,tokenMaker,config)
	require.NoError(t, err)
	return handler
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())

}