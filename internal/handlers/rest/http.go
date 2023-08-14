package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/internal/core/ports"
	"github.com/kamalbowselvam/chaintask/util"
)

type HttpHandler struct {
	taskService ports.TaskService
	userService ports.UserService
}

func NewHttpHandler(taskService ports.TaskService, userService ports.UserService) *HttpHandler {

	return &HttpHandler{
		taskService: taskService,
		userService: userService,
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

type CreateUserParams struct {
	Name     string `json:"name"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
	Email    string `json:"mail"`
}

func (h *HttpHandler) CreateUser(c *gin.Context) {
	taskparam := CreateUserParams{}
	c.BindJSON(&taskparam)
	hashed_pass, err := util.HashPassword(taskparam.Password)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Could not hash password"})
	}
	user, err := h.userService.CreateUser(taskparam.Name, hashed_pass, taskparam.FullName, taskparam.Email)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
	}
	c.JSON(200, user)
}

type LoginParams struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (h *HttpHandler) Login(c *gin.Context) {
	loginparams := LoginParams{}
	c.BindJSON(&loginparams)
	user, err := h.userService.GetUser(loginparams.UserName)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
	}
	// FIXME Token Creation Here
	if util.CheckPassword(loginparams.Password, user.HashedPassword) == nil {
		token := "token_to_be_created"
		c.JSON(200, token)
	} else {
		c.AbortWithStatusJSON(403, gin.H{"message": "Unauthorized"})
	}

}
