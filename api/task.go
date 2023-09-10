package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
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
// @Router       /tasks/{id} [get]
// @Security BearerAuth
func (s *Server) GetTask(c *gin.Context) {
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
// @Router       /tasks/{id} [delete]
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
// @Router       /tasks/ [post]
// @Security BearerAuth
func (s *Server) CreateTask(c *gin.Context) {
	taskparam := db.CreateTaskParams{}
	c.ShouldBindBodyWith(&taskparam, binding.JSON)

	task, err := s.service.CreateTask(c, taskparam)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, task)
}

func (s *Server) UpdateTask(c *gin.Context) {
	taskparam := domain.Task{}
	c.BindJSON(&taskparam)
	log.Println(taskparam)

	task, err := s.service.UpdateTask(c, taskparam)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, task)
}
