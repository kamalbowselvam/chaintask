package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
)

// CreateProject godoc
// @Summary      Create a project
// @Description  create a  project
// @Tags         projects
// @Produce      json
// @Param        companyId path     int true   "Company ID"
// @Param        request body db.CreateProjectParam true "project creation parameters"
// @Success      200  {object}  domain.Project
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /company/{companyId}/projects/ [post]
// @Security BearerAuth
func (s *Server) CreateProject(c *gin.Context) {

	projectparam := db.CreateProjectParam{}

	err := c.ShouldBindBodyWith(&projectparam, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var req db.CompanyParam
	err = c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	createdBy, existed := c.Get(authorizationPayloadKey)
	if !existed {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Forbidden": ""})
	}

	projectparam.CompanyId = req.CompanyId
	projectparam.CreatedBy = createdBy.(*token.Payload).Username;
	task, err := s.service.CreateProject(c, projectparam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, task)
}
