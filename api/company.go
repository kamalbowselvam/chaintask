package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/logger"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
)

// CreateCompany godoc
// @Summary      Create a company
// @Description  Creates a company
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        request body domain.Company true "company creation parameter"
// @Success      200  {object}  domain.Company
// @Failure      400  {object} error
// @Failure      404  {object} error
// @Failure      500  {object} error
// @Router       /companies [post]
// @Security BearerAuth
func (s *Server) CreateCompany(c *gin.Context){

	companyparams := db.CreateCompanyParams{}
	err := c.BindJSON(&companyparams)
	logger.Info(companyparams.CompanyName)
	
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	createdBy, existed := c.Get(authorizationPayloadKey)
	if !existed {
		c.AbortWithStatusJSON(http.StatusForbidden, util.ErrorResponseString("Forbidden"))
	}
	companyparams.CreatedBy = createdBy.(*token.Payload).Username;
	company, err := s.service.CreateCompany(c, companyparams)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	c.JSON(200, company)
}