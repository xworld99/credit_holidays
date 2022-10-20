package main

import (
	"credit_holidays/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	log "github.com/sirupsen/logrus"
	"os"
)

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
	if err != nil {
		log.WithError(err).WithField("path", confPath).Fatal("can't init handler")
	}

	r := gin.Default()

	userGroup := r.Group("/user")
	userGroup.GET("/get_balance", handler.GetBalance)

	orderGroup := r.Group("/order")
	orderGroup.POST("/add_order", handler.AddOrder)
	orderGroup.POST("/change_order_status", handler.ChangeOrderStatus)

	r.Static("/reports", cfg.String("path.static"))

	log.Fatal(r.Run(":8080"))
}
