package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/logger"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"go.uber.org/zap"
)

// GetTask godoc
// @Summary      Get a Task by its ID
// @Description  get a task by its ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  domain.Task
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /company/{companyID}/projects/{projectId}/tasks/{taskId} [get]
// @Security BearerAuth
func (s *Server) GetTask(c *gin.Context) {
	logger := logger.FromCtx(c.Request.Context())
	logger.Info("toto")
	var req db.GetTaskParams
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	task, err := s.service.GetTask(c, req.Id)

	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return

	}

	authorizationPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if task.CreatedBy != authorizationPayload.Username {
		err := errors.New("task does not belong to user")
		c.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary      Delete a Task by its ID
// @Description  delete a task by its ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      202
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /company/{companyID}/projects/{projectId}/tasks/{id} [delete]
// @Security BearerAuth
func (s *Server) DeleteTask(c *gin.Context) {
	var req db.GetTaskParams
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	err = s.service.DeleteTask(c, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusAccepted, nil)
}

// CreateTask godoc
// @Summary      Create a Task
// @Description  Create a tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body db.CreateTaskParams true "task creation parameter"
// @Success      200  {object}  domain.Task
// @Failure      400  {object} error
// @Failure      404  {object} error
// @Failure      500  {object} error
// @Router       /company/{companyID}/projects/{projectId}/tasks/ [post]
// @Security BearerAuth
func (s *Server) CreateTask(c *gin.Context) {
	taskparam := db.CreateTaskParams{}
	c.ShouldBindBodyWith(&taskparam, binding.JSON)

	createdBy, existed := c.Get(authorizationPayloadKey)
	if !existed {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Forbidden": ""})
	}
	taskparam.CreatedBy = createdBy.(*token.Payload).Username;

	logger.Debug("Creating a task",
		zap.String("package", "api"),
		zap.String("function", "CreateTask"),
		zap.Any("param",taskparam),
	)

	task, err := s.service.CreateTask(c, taskparam)

	if err != nil {

		logger.Error("Creating a task",
			zap.String("package", "api"),
			zap.String("function", "CreateTask"),
			zap.Any("param",taskparam),
		)
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, task)
}

// UpdateTask godoc
// @Summary      Update a Task
// @Description  Updates a tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body domain.Task true "task update parameter"
// @Success      200  {object}  domain.Task
// @Failure      400  {object} error
// @Failure      404  {object} error
// @Failure      500  {object} error
// @Router       /company/{companyId}/projects/{projectId}/tasks/{taskId} [post]
// @Security BearerAuth
func (s *Server) UpdateTask(c *gin.Context) {
	taskparam := db.UpdateTaskParams{}
	c.BindJSON(&taskparam)
	//s.logger.Sugar().Info(taskparam)

	task, err := s.service.UpdateTask(c, taskparam)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, task)
}
