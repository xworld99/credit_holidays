package db

import (
	"context"
	"credit_holidays/internal/models"
	"database/sql"
	"time"
)

type CreditHolidaysDB interface {
	// Begin initiates new transaction with specific isolation level
	Begin(context.Context, sql.IsolationLevel) (*sql.Tx, error)
	// Commit commits passed transaction
	Commit(*sql.Tx)
	// Rollback rollbacks passed transaction
	Rollback(*sql.Tx)
	// Close closes db connection
	Close()

	// GetUserById extract user with specific id from db
	GetUserById(context.Context, models.User) (models.User, error)
	// UpdateUser save changes of user`s balance to db, ALLOWED ONLY IN TRANSACTION
	UpdateUser(context.Context, *sql.Tx, models.User) (models.User, error)
	// InsertUserIfNotExists save new user with specific id in db, ALLOWED ONLY IN TRANSACTION
	InsertUserIfNotExists(context.Context, *sql.Tx, models.User) (models.User, error)

	// GetServicesList extract information about all services from db
	GetServicesList(context.Context) ([]models.Service, error)
	// GetServiceById extract information about specific service from db
	GetServiceById(context.Context, models.Service) (models.Service, error)

	// CreateOrder saves information about new order in db, ALLOWED ONLY IN TRANSACTION
	CreateOrder(context.Context, *sql.Tx, models.Order) (models.Order, error)
	// UpdateOrder saves updated information about order in db, ALLOWED ONLY IN TRANSACTION
	UpdateOrder(context.Context, *sql.Tx, models.Order) (models.Order, error)
	// GetFullOrderInfo extract whole information about order (order, user, service) from db,
	//   ALLOWED ONLY IN TRANSACTION
	GetFullOrderInfo(
		context.Context,
		*sql.Tx,
		models.Order,
		models.User,
		models.Service,
	) (models.Order, models.User, models.Service, error)
	// GetLastOrderMonth extract last order of specific moth from db, used in generating csv reports
	GetLastOrderMonth(context.Context, time.Time) (models.Order, error)

	// FormReport extract information about services from db
	FormReport(context.Context, models.CSVData) (models.CSVData, error)
	// GetHistoryFrame extract information about user history from db
	GetHistoryFrame(context.Context, models.HistoryFrame) (models.HistoryFrame, error)
}
