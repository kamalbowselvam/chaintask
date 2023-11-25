package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/logger"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "author"
	requestUUIDKey          = "requestUUID"
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
func AuthorizeMiddleware(authorize authorization.AuthorizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		val, existed := c.Get(authorizationPayloadKey)
		if !existed {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponseString("user has not logged in yet"))
			return
		}
		ok, err := authorize.Enforce(val.(*token.Payload), c.Request.URL.Path, c.Request.Method)
		if err != nil {
			logger.Warn("Error occured while authorizing the user ", zap.Error(err))
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


// add logger middleware
// Inspired from https://betterstack.com/community/guides/logging/go/zap/#adding-context-to-your-logs
// and https://github.com/betterstack-community/go-logging/blob/zap/middleware.go

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func requestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // retrieve the standard logger instance
        l := logger.Get()

        // create a correlation ID for the request
        correlationID := xid.New().String()

        c.Set(requestUUIDKey, correlationID)

        // create a child logger containing the correlation ID
        // so that it appears in all subsequent logs
        l = l.With(zap.String(string(requestUUIDKey), correlationID))

        c.Writer.Header().Add("X-Correlation-ID", correlationID)

        lrw := newLoggingResponseWriter(c.Writer)

        // the logger is associated with the request context here
        // so that it may be retrieved in subsequent `http.Handlers`

		r := c.Request
        
		defer func(start time.Time) {
			l.Info(
				fmt.Sprintf(
					"%s request to %s completed",
					r.Method,
					r.RequestURI,
				),
				zap.String("method", r.Method),
				zap.String("url", r.RequestURI),
				zap.String("user_agent", r.UserAgent()),
				zap.Int("status_code", lrw.statusCode),
				zap.Duration("elapsed_ms", time.Since(start)),
			)
		}(time.Now())


        c.Next()
    }
}
