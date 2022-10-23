package handlers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GenerateReport godoc
// @Summary generate report
// @Schemes
// @Description generate report with description of services in specific month in format "MM-YYYY"
// @Tags report
// @Accept json
// @Produce json
// @Param month query string true "month of year in format MM-YYYY"
// @Success 200 {string} string "path to generated report in static dir"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} sting  "internal server error"
// @Router /report/generate_report [get]
func (h *Handler) GenerateReport(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	month, ok := params["month"]
	if !ok {
		log.WithError(fmt.Errorf("cant find month in req query")).Error("cant process request body")
		ctx.JSON(http.StatusBadRequest, consts.ErrorDescriptions[http.StatusBadRequest])
		return
	}

	resp, err := h.ctrl.GenerateReport(ctx, models.GenerateReportRequest(month[0]))
	if err.Err != nil {
		log.WithError(err.Err).Error("cant generate report")
		ctx.JSON(err.Type, consts.ErrorDescriptions[err.Type])
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
