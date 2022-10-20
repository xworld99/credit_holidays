package handlers

import (
	"credit_holidays/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetBalance(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	data, ok := params["id"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract user id from query")
		return
	}

	u := models.GetBalanceRequest(data[0])
	resp, err := h.ctrl.GetBalance(ctx, u)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
