package handlers

import (
	"credit_holidays/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) AddOrder(ctx *gin.Context) {
	var reqData models.AddOrderRequest

	if err := ctx.BindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid request body")
		return
	}

	// create order
	resp, err := h.ctrl.AddOrder(ctx, reqData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("%w", err))
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) ChangeOrderStatus(ctx *gin.Context) {
	var reqData models.ChangeOrderRequest

	if err := ctx.BindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid request body")
		return
	}

	// change order status
	resp, err := h.ctrl.ChangeOrderStatus(ctx, reqData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
