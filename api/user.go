package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/logger"
	"github.com/kamalbowselvam/chaintask/util"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type createUserRequest struct {
	Username  string `json:"user_name" binding:"required,alphanum"`
	Password  string `json:"password" binding:"required,min=6"`
	FullName  string `json:"full_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	UserRole  string `json:"user_role" binding:"required,user_role"`
	CompanyId *int64 `json:"company_id" binding:"required,number"`
}

type userResponse struct {
	Username          string    `json:"user_name"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	UserRole          string    `json:"user_role"`
	CompanyId         int64     `json:"company_id"`
}

func newUserResponse(user domain.User) userResponse {

	rsp := userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		UserRole:          user.UserRole,
		CompanyId:         user.CompanyId,
	}

	return rsp
}

// CreateUser godoc
// @Summary      Create a User
// @Description  Create a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body createUserRequest true "user creation parameters"
// @Success      200  {object}  domain.User
// @Failure      400  {object} error
// @Failure      404  {object} error
// @Failure      500  {object} error
// @Router       /users/ [post]
func (s *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		UserRole:       req.UserRole,
		CompanyId:      *req.CompanyId,
	}

	user, err := s.service.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

// LoginUser godoc
// @Summary      Logs an user in
// @Description  Logs an user in
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body loginUserRequest true "user creation parameters"
// @Success      200  {object}  loginUserResponse
// @Failure      400  {object} error
// @Failure      404  {object} error
// @Failure      500  {object} error
// @Router       /users/login [post]
func (s *Server) LoginUser(ctx *gin.Context) {

	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	logger := logger.FromGinCtx(ctx)
	logger.Debug("Login", zap.Any("login", req))

	user, err := s.service.GetUser(ctx, req.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))

	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		user.UserRole,
		s.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		user.UserRole,
		s.config.RefreshTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	session, err := s.service.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)

}
