package controllers

import (
	"credit_holidays/internal/db"
	"github.com/knadh/koanf"
	"time"
)

// Controller provides all business logic for handler struct
type Controller struct {
	db            db.CreditHolidaysDB
	selectTm      time.Duration
	updateTm      time.Duration
	deleteTm      time.Duration
	insertTm      time.Duration
	transactionTm time.Duration
}

func NewService(cfg *koanf.Koanf, db db.CreditHolidaysDB) (*Controller, error) {
	return &Controller{
		db:            db,
		selectTm:      cfg.Duration("timeout.db_select"),
		updateTm:      cfg.Duration("timeout.db_update"),
		deleteTm:      cfg.Duration("timeout.db_delete"),
		insertTm:      cfg.Duration("timeout.db_insert"),
		transactionTm: cfg.Duration("timeout.db_transaction"),
	}, nil
}
