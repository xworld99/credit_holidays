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
	// AddOrder function save new order in db
	AddOrder(context.Context, models.AddOrderRequest) (models.Order, models.HandlerError)
	// ChangeOrderStatus allows to change status of order from "in_progress" to "success" or "declined"
	ChangeOrderStatus(context.Context, models.ChangeOrderRequest) (models.Order, models.HandlerError)

	// GenerateReport generates report in csv format and return link on this report
	GenerateReport(context.Context, models.GenerateReportRequest) (string, models.HandlerError)

	// GetServicesList return all services from db
	GetServicesList(context.Context) ([]models.Service, models.HandlerError)
	// GetServiceInfo extract information about specific service, modifies passed service
	GetServiceInfo(context.Context, *models.Service) error

	// GetBalance return balance, frozen_balance of certain user
	GetBalance(context.Context, models.GetBalanceRequest) (models.User, models.HandlerError)
	// GetHistory return history of users operations
	GetHistory(context.Context, models.GetHistoryRequest) (models.HistoryFrame, models.HandlerError)
	// InsertUserIfNotExists creates new user in db if it isn`t exists, modifies passed user
	InsertUserIfNotExists(context.Context, *sql.Tx, *models.User) error

	// Close method closes connection to db if it exists
	Close()
}

// Controller provides all business logic for handler struct
type Controller struct {
	db         db.CreditHolidaysDB
	dbTm       time.Duration
	staticPath string
}

// NewController simple constructor for Controller, used in wire code generation
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
