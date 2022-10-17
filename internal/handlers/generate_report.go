package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GenerateReport(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, struct{}{})
}
