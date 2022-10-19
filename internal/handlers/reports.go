package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GenerateReport(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, struct{}{})
}

func (h *Handler) GetUserHistory(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, struct{}{})
}

func (h *Handler) GetServicesList(ctx *gin.Context) {
	resp, err := h.ctrl.GetServicesList(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("%w", err))
	}

	ctx.JSON(http.StatusOK, resp)
}
