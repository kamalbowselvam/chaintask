package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/service"
)

type HttpHandler struct {
	taskService service.TaskService
}

func NewHttpHandler(taskService service.TaskService) *HttpHandler {

	return &HttpHandler{
		taskService: taskService,
	}
}

type getTaskRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (h *HttpHandler) GetTask(c *gin.Context) {

	var req getTaskRequest

	err := c.ShouldBindUri(&req)
	// id, err := strconv.ParseInt(c.Param("id"),10,64)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	task, err := h.taskService.GetTask(id)
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, task)
}




func (h *HttpHandler) CreateTask(c *gin.Context) {
	taskparam := db.CreateTaskParams{}
	c.BindJSON(&taskparam)

	task, err := h.taskService.CreateTask(c,taskparam)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, task)
}
