package api

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/db"
)





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