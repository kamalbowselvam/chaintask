package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/token"
)

// CreateProject godoc
// @Summary      Create a project
// @Description  create a  project
// @Tags         projects
// @Produce      json
// @Param        request body db.CreateProjectParam true "project creation parameters"
// @Success      200  {object}  domain.Project
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /projects/ [post]
// @Security BearerAuth
func (s *Server) CreateProject(c *gin.Context) {

	projectparam := db.CreateProjectParam{}
	c.ShouldBindBodyWith(&projectparam, binding.JSON)
	log.Println(projectparam)

	createdBy, existed := c.Get(authorizationPayloadKey)
	if !existed {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Forbidden": ""})
	}
	projectparam.CreatedBy = createdBy.(*token.Payload).Username;
	task, err := s.service.CreateProject(c, projectparam)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, task)
}
