// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package handlers

import (
	"credit_holidays/internal/controllers"
	"credit_holidays/internal/db"
	"github.com/knadh/koanf"
)

// Injectors from wire.go:

func InitializeHandler(cfg *koanf.Koanf) (*Handler, error) {
	creditHolidaysDB, err := db.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}
	creditHolidaysController := controllers.NewController(cfg, creditHolidaysDB)
	handler := NewHandler(creditHolidaysController)
	return handler, nil
}
