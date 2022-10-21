package handlers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

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
