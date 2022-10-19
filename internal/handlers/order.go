package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) AddOrder(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	// parsing query
	userId, ok := params["user_id"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract user id from query")
	}

	serviceId, ok := params["service_id"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract service id from query")
	}

	amount, ok := params["amount"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract amount from query")
	}

	// create order
	resp, err := h.ctrl.AddOrder(ctx, userId[0], serviceId[0], amount[0])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("%w", err))
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) ChangeOrderStatus(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	// parsing query
	orderId, ok := params["order_id"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract order id from query")
	}

	action, ok := params["action"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract order id from query")
	}

	// create order
	resp, err := h.ctrl.ChangeOrderStatus(ctx, orderId[0], action[0])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("%w", err))
	}

	ctx.JSON(http.StatusOK, resp)
}
