package main

import (
	_ "credit_holidays/api"
	"credit_holidays/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// @title CreditHolidaysAPI
// @version 1.0
// @description API for Credit Holidays app
// @host 0.0.0.0:8080
// @basePath /
func main() {
	confPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		log.Fatal("can't get config path from environment")
	}

	cfg := koanf.New(".")
	if err := cfg.Load(file.Provider(confPath), yaml.Parser()); err != nil {
		log.WithError(err).WithField("path", confPath).Fatal("can't read config")
	}

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	})

	if level, err := log.ParseLevel(cfg.String("log.level")); err != nil {
		log.WithError(err).Fatal("can't get log level from config")
	} else {
		log.SetLevel(level)
	}

	handler, err := handlers.InitializeHandler(cfg)
	defer handler.Close()

	if err != nil {
		log.WithError(err).WithField("path", confPath).Fatal("can't init handler")
	}

	r := gin.Default()

	r.Use(CORSMiddleware())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/reports", cfg.String("path.static"))

	userGroup := r.Group("/user")
	userGroup.GET("/get_balance", handler.GetBalance)
	userGroup.GET("/get_history", handler.GetUserHistory)

	orderGroup := r.Group("/order")
	orderGroup.POST("/add_order", handler.AddOrder)
	orderGroup.POST("/change_order_status", handler.ChangeOrderStatus)

	reportGroup := r.Group("/report")
	reportGroup.GET("/generate_report", handler.GenerateReport)

	serviceGroup := r.Group("/service")
	serviceGroup.GET("/get_all", handler.GetServicesList)

	log.Fatal(r.Run(":8080"))
}
