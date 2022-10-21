package controllers

import (
	"credit_holidays/internal/db"
	"github.com/knadh/koanf"
	"time"
)

// Controller provides all business logic for handler struct
type Controller struct {
	db         db.CreditHolidaysDB
	dbTm       time.Duration
	staticPath string
}

func NewController(cfg *koanf.Koanf, db db.CreditHolidaysDB) (*Controller, error) {
	return &Controller{
		db:         db,
		dbTm:       cfg.Duration("timeout.db_transaction"),
		staticPath: cfg.String("path.static"),
	}, nil
}

func (c *Controller) Close() {
	if c.db != nil {
		c.db.Close()
	}
}
