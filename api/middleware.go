package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "author"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// Authorize determines if current subject has been authorized to take an action on an object.
func AuthorizeMiddleware(obj interface{}, act string, adapter persist.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		val, existed := c.Get(authorizationPayloadKey)
		if !existed {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponseString("user has not logged in yet"))
			return
		}
		// Casbin enforces policy
		log.Println(val)
		err := c.BindJSON(&obj)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.ErrorResponse(err))
			return
		}
		ok, err := enforce(val.(*token.Payload), obj, act, adapter)
		if err != nil {
			log.Fatalf("Error occured while authorizing the user %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.ErrorResponse(err))
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, util.ErrorResponseString("forbidden"))
			return
		}
		c.Next()
	}
}

func enforce(sub *token.Payload, obj interface{}, act string, adapter persist.Adapter) (bool, error) {
	// Load model configuration file and policy store adapter

	enforcer, err := casbin.NewEnforcer("./config/rbac_model.conf", adapter)
	if err != nil {
		log.Fatal(err)
	}
	// Load policies from DB dynamically
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	// Verify
	ok, err := enforcer.Enforce(sub.Role, obj, act)
	return ok, err
}
