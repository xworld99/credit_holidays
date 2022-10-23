package controllers

import (
	"context"
	"credit_holidays/internal/db"
	"credit_holidays/internal/models"
	"database/sql"
	"github.com/knadh/koanf"
	"time"
)

type CreditHolidaysController interface {
	// orders
	AddOrder(context.Context, models.AddOrderRequest) (models.Order, models.HandlerError)
	ChangeOrderStatus(context.Context, models.ChangeOrderRequest) (models.Order, models.HandlerError)

	// reports
	GenerateReport(context.Context, models.GenerateReportRequest) (string, models.HandlerError)

	// services
	GetServicesList(context.Context) ([]models.Service, models.HandlerError)
	GetServiceInfo(context.Context, *models.Service) error

	// user
	GetBalance(context.Context, models.GetBalanceRequest) (models.User, models.HandlerError)
	GetHistory(context.Context, models.GetHistoryRequest) (models.HistoryFrame, models.HandlerError)
	InsertUserIfNotExists(context.Context, *sql.Tx, *models.User) error

	// general
	Close()
}

// Controller provides all business logic for handler struct
type Controller struct {
	db         db.CreditHolidaysDB
	dbTm       time.Duration
	staticPath string
}

func NewController(cfg *koanf.Koanf, db db.CreditHolidaysDB) CreditHolidaysController {
	return &Controller{
		db:         db,
		dbTm:       cfg.Duration("timeout.db_transaction"),
		staticPath: cfg.String("path.static"),
	}
}

func (c *Controller) Close() {
	if c.db != nil {
		c.db.Close()
	}
}
