package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
)

type HttpHandler struct {
	taskService service.TaskService
	tokenMaker  token.Maker
	config      util.Config
}

func NewHttpHandler(taskService service.TaskService, tokenMaker token.Maker, config util.Config) *HttpHandler {

	return &HttpHandler{
		taskService: taskService,
		tokenMaker:  tokenMaker,
		config:      config,
	}
}

func (h *HttpHandler) GetTokenMaker() *token.Maker {
	return &h.tokenMaker
}

func (h *HttpHandler) GetTask(c *gin.Context) {
	var req db.GetTaskParams
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	task, err := h.taskService.GetTask(c, req.Id)

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

func (h *HttpHandler) CreateTask(c *gin.Context) {
	taskparam := db.CreateTaskParams{}
	c.BindJSON(&taskparam)
	log.Println(taskparam)

	task, err := h.taskService.CreateTask(c, taskparam)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, task)
}


func (h *HttpHandler) CreateProject(c *gin.Context){
	projectparam := db.CreateProjectParam{}
	c.BindJSON(&projectparam)
	log.Println(projectparam)

	task, err := h.taskService.CreateProject(c, projectparam)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, task)
}