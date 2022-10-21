package handlers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) AddOrder(ctx *gin.Context) {
	var reqData models.AddOrderRequest

	if err := ctx.BindJSON(&reqData); err != nil {
		log.WithError(err).Error("cant process request body")
		ctx.JSON(http.StatusBadRequest, consts.ErrorDescriptions[http.StatusBadRequest])
		return
	}

	// create order
	resp, err := h.ctrl.AddOrder(ctx, reqData)
	if err.Err != nil {
		log.WithError(err.Err).Error("cant create order")
		ctx.JSON(err.Type, consts.ErrorDescriptions[err.Type])
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) ChangeOrderStatus(ctx *gin.Context) {
	var reqData models.ChangeOrderRequest

	if err := ctx.BindJSON(&reqData); err != nil {
		log.WithError(err).Error("cant process request body")
		ctx.JSON(http.StatusBadRequest, consts.ErrorDescriptions[http.StatusBadRequest])
		return
	}

	// change order status
	resp, err := h.ctrl.ChangeOrderStatus(ctx, reqData)
	if err.Err != nil {
		log.WithError(err.Err).Error("cant change order status")
		ctx.JSON(err.Type, consts.ErrorDescriptions[err.Type])
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
