package handlers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// AddOrder godoc
// @Summary create new order for user with specific service
// @Schemes
// @Description initiates a change in the user's balance, returns order info
// @Tags order
// @Accept json
// @Produce json
// @Param info body models.AddOrderRequest true "Info about order"
// @Success 200 {object} models.Order "created order"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} sting  "internal server error"
// @Router /order/add_order [post]
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

// ChangeOrderStatus godoc
// @Summary change status of existing order
// @Schemes
// @Description proof or decline existing order, return current state of order
// @Tags order
// @Accept json
// @Produce json
// @Param info body models.ChangeOrderRequest true "Info about order"
// @Success 200 {object} models.Order "proofed or declined order"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} sting  "internal server error"
// @Router /order/change_order_status [post]
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
