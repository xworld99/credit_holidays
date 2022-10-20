package handlers

import (
	"credit_holidays/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// month-year -> csv file
func (h *Handler) GenerateReport(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, struct{}{})
}

func (h *Handler) GetUserHistory(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	userId, ok := params["user_id"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract user id from query")
		return
	}
	fromDate, ok := params["from_date"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract from date from query")
		return
	}
	toDate, ok := params["to_date"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "cant extract to date from query")
		return
	}
	orderBy, ok := params["order_by"]
	if !ok {
		orderBy = []string{"created_at"}
	}
	offset, ok := params["offset"]
	if !ok {
		offset = []string{"0"}
	}
	limit, ok := params["limit"]
	if !ok {
		limit = []string{"10"}
	}

	resp, err := h.ctrl.GetHistory(ctx, models.GetHistoryRequest{
		UserId:   userId[0],
		FromDate: fromDate[0],
		ToDate:   toDate[0],
		Offset:   offset[0],
		Limit:    limit[0],
		OrderBy:  orderBy[0],
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetServicesList(ctx *gin.Context) {
	resp, err := h.ctrl.GetServicesList(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("%w", err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
