package main

import (
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

	/*
		handler = handlers.InitializeHandler(cfg)

		r := gin.Default()

		userGroup := r.Group("/user")
		userGroup.GET("/get_balance", handler.GetBalance)
		userGroup.GET("/get_history", handler.GetHistory)
		userGroup.POST("/add_money", handler.AddMoney)
		userGroup.POST("/withdraw_money", handler.WithdrawMoney)
		userGroup.POST("/transfer_money", handler.TransferMoney)

		orderGroup := r.Group("/order")
		orderGroup.POST("/init_order", handler.InitOrder)
		orderGroup.POST("/change_order_status", handler.ChangeOrderStatus)

		utilsGroup := r.Group("/utils")
		utilsGroup.GET("/generate_report", handler.GenerateReport)

		log.Fatal(r.Run(":8080"))
	*/
}
