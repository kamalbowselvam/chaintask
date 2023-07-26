package rest

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/internal/core/ports"
)

type HttpHandler struct {
	taskService ports.TaskService
}

func NewHttpHandler(taskService ports.TaskService) *HttpHandler {

	return &HttpHandler{
		taskService: taskService,
	}
}


type getTaskRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}


func(h *HttpHandler) GetTask(c *gin.Context){
	
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

type TaskParams struct {
	Name      string  `json:"name"`
	Budget    float64 `json:"budget"`
	CreatedBy string  `json:"createdBy"`
}

func (h *HttpHandler) CreateTask(c *gin.Context) {
	taskparam := TaskParams{}
	c.BindJSON(&taskparam)

	task, err := h.taskService.CreateTask(taskparam.Name, taskparam.Budget, taskparam.CreatedBy)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, task)
}
