package handlers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetBalance godoc
// @Summary access user balance
// @Schemes
// @Description return balances of user if it exists
// @Tags user
// @Accept json
// @Produce json
// @Param id query int true "id of user"
// @Success 200 {object} models.User "info about user"
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} sting  "internal server error"
// @Router /user/get_balance [get]
func (h *Handler) GetBalance(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	data, ok := params["id"]
	if !ok {
		log.WithError(fmt.Errorf("cant find id in req query")).Error("cant parse user request")
		ctx.JSON(http.StatusBadRequest, consts.ErrorDescriptions[http.StatusBadRequest])
		return
	}

	u := models.GetBalanceRequest(data[0])
	resp, err := h.ctrl.GetBalance(ctx, u)
	if err.Err != nil {
		log.WithError(err.Err).Error("cant get users balance")
		ctx.JSON(err.Type, consts.ErrorDescriptions[err.Type])
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetUserHistory godoc
// @Summary return history of user`s orders
// @Schemes
// @Description return list of orders attached to specific user
// @Tags user
// @Accept json
// @Produce json
// @Param id query int true "id of user"
// @Param from_date query string true "start date in format DD-MM-YYYY"
// @Param to_date query string true "end date in format DD-MM-YYYY"
// @Param order_by query string false "sorting type for orders = amount, created_at, default = created_at"
// @Param limit query int false "max orders in response default 10"
// @Param offset query int false "offset default 0"
// @Success 200 {object} models.HistoryFrame "info about user"
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} sting  "internal server error"
// @Router /user/get_history [get]
func (h *Handler) GetUserHistory(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	userId, ok := params["id"]
	if !ok {
		log.WithError(fmt.Errorf("cant find id in req query")).Error("cant process request body")
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
		log.WithError(err.Err).Error("cant extract users history")
		ctx.JSON(err.Type, consts.ErrorDescriptions[err.Type])
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
