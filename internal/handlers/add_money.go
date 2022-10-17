package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddMoney(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, struct{}{})
}
