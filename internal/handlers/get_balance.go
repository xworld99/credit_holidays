package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBalance(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, struct{}{})
}
