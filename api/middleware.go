package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "author"
)

type WriteDetail struct {
	CreatedBy string
}

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
func AuthorizeMiddleware(act string, adapter interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovering from panic")
			}
		}()
		// Get current user/subject
		val, existed := c.Get(authorizationPayloadKey)
		if !existed {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponseString("user has not logged in yet"))
			return
		}
		// Casbin enforces policy
		var obj interface{}
		c.ShouldBindBodyWith(&obj, binding.JSON)
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

func enforce(sub *token.Payload, abstract interface{}, act string, adapter interface{}) (bool, error) {
	// Load model configuration file and policy store adapter
	conf_file_path := "./config/rbac_model.conf"
	fmt.Println(sub)
	_, err := os.Stat(conf_file_path)
	if err != nil {
		conf_file_path = "." + conf_file_path
	}
	enforcer, err := casbin.NewEnforcer(conf_file_path, adapter)
	if err != nil {
		log.Fatal(err)
	}
	// Load policies from DB dynamically
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	// Verify
	temp := abstract.(map[string]interface{})
	res, check := temp["CreatedBy"].(string)
	obj := WriteDetail{}
	if check {
		obj.CreatedBy = res
	} else {
		obj.CreatedBy = sub.Username
	}
	ok, err := enforcer.Enforce(sub, obj, act)
	return ok, err
}
