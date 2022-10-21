package handlers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) GenerateReport(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, struct{}{})
}

func (h *Handler) GetUserHistory(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	userId, ok := params["user_id"]
	if !ok {
		log.WithError(fmt.Errorf("cant find user_id in req query")).Error("cant process request body")
		ctx.JSON(http.StatusBadRequest, consts.ErrorDescriptions[http.StatusBadRequest])
		return
	}
	fromDate, ok := params["from_date"]
	if !ok {
		log.WithError(fmt.Errorf("cant find from_date in req query")).Error("cant process request body")
		ctx.JSON(http.StatusBadRequest, consts.ErrorDescriptions[http.StatusBadRequest])
		return
	}
	toDate, ok := params["to_date"]
	if !ok {
		log.WithError(fmt.Errorf("cant find to_date in req query")).Error("cant process request body")
		ctx.JSON(http.StatusBadRequest, consts.ErrorDescriptions[http.StatusBadRequest])
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

	if err.Err != nil {
		ctx.JSON(err.Type, consts.ErrorDescriptions[err.Type])
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetServicesList(ctx *gin.Context) {
	resp, err := h.ctrl.GetServicesList(ctx)
	if err.Err != nil {
		log.WithError(err.Err).Error("cant extract services list from db")
		ctx.JSON(err.Type, consts.ErrorDescriptions[err.Type])
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
