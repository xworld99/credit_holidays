package main

import (
	"credit_holidays/internal/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	r.GET("/get_balance", handlers.GetBalance)
	r.GET("/get_history", handlers.GetHistory)
	r.POST("/add_money", handlers.AddMoney)

	r.POST("/init_order", handlers.InitOrder)
	r.POST("/proof_order", handlers.ProofOrder)

	r.GET("/generate_repost", handlers.GenerateReport)

	log.Fatal(r.Run(":8080"))
}
