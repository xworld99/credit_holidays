package handlers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetServicesList godoc
// @Summary get info about services
// @Schemes
// @Description return list of all available services
// @Tags service
// @Accept json
// @Produce json
// @Success 200 {array} models.Service "info about services"
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} sting  "internal server error"
// @Router /service/get_all [get]
func (h *Handler) GetServicesList(ctx *gin.Context) {
	var resp []models.Service

	resp, err := h.ctrl.GetServicesList(ctx)
	if err.Err != nil {
		log.WithError(err.Err).Error("cant extract services list from db")
		ctx.JSON(err.Type, consts.ErrorDescriptions[err.Type])
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
