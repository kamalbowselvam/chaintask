package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/logger"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	"go.uber.org/zap"
)

// Createpayment godoc
// @Summary      Pay for a task
// @Description  Adds a payment source and a payment amout for a task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        taskId   path      int  true  "Task ID"
// @Param        projectId path     int true   "Project ID"
// @Param        companyId path     int true   "Company ID"
// @Param        request body db.CreateTaskPaymentParams true "payment source"
// @Success      200  {object}  domain.TaskPayment
// @Failure      400  {object} error
// @Failure      404  {object} error
// @Failure      500  {object} error
// @Router       /payment/{paymentId}/projects/{projectId}/payment/{taskId} [post]
// @Security BearerAuth
func (s *Server) PayForATask(c *gin.Context){

	paymentparams := db.CreateTaskPaymentParams{}
	err := c.BindJSON(&paymentparams)
	logger.Info("Task payment api", zap.Any("task payment params", paymentparams));
	
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	createdBy, existed := c.Get(authorizationPayloadKey)
	if !existed {
		c.AbortWithStatusJSON(http.StatusForbidden, util.ErrorResponseString("Forbidden"))
	}
	paymentparams.CreatedBy = createdBy.(*token.Payload).Username;
	payment, err := s.service.PayForATask(c, paymentparams)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	c.JSON(200, payment)
}