package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetBalance(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	userId, ok := params["id"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract user id from query")
	}

	resp, err := h.ctrl.GetBalance(ctx, userId[0])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	ctx.JSON(http.StatusOK, resp)
}
