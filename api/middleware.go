package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "author"
)

type IdObject struct {
	Id int64 `json:"Id" uri:"id"`
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
func AuthorizeMiddleware(authorize authorization.AuthorizationService) gin.HandlerFunc {
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
		ok, err := authorize.Enforce(val.(token.Payload), c.FullPath(), c.Request.Method)
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
