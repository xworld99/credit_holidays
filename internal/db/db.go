package db

import (
	"context"
	"credit_holidays/internal/models"
	"database/sql"
	"time"
)

type CreditHolidaysDB interface {
	// general
	Begin(context.Context, sql.IsolationLevel) (*sql.Tx, error)
	CommitTransaction(*sql.Tx) error
	RollbackTransaction(*sql.Tx) error

	// users interaction
	GetUserById(context.Context, models.User) (models.User, error)
	UpdateUser(context.Context, *sql.Tx, models.User) (models.User, error)
	GetCreateUser(context.Context, *sql.Tx, models.User) (models.User, error)

	// services interaction
	GetServicesList(context.Context) ([]models.Service, error)
	GetServiceById(context.Context, models.Service) (models.Service, error)

	// orders interaction
	CreateOrder(context.Context, *sql.Tx, models.Order) (models.Order, error)
	UpdateOrder(context.Context, *sql.Tx, models.Order) (models.Order, error)
	GetFullOrderInfo(
		context.Context,
		*sql.Tx,
		models.Order,
		models.User,
		models.Service,
	) (models.Order, models.User, models.Service, error)
	GetLastOrderMonth(context.Context, time.Time) (models.Order, error)

	// reports forming
	FormReport(context.Context, models.CSVData) (models.CSVData, error)
	GetHistoryFrame(context.Context, models.HistoryFrame) (models.HistoryFrame, error)
}
