package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitOrder(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, struct{}{})
}
