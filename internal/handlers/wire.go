//go:build wireinject
// +build wireinject

package handlers

import (
	"credit_holidays/internal/controllers"
	"credit_holidays/internal/db"
	"github.com/google/wire"
	"github.com/knadh/koanf"
)

func InitializeHandler(cfg *koanf.Koanf) (*Handler, error) {
	wire.Build(
		db.NewPostgresDB,
		controllers.NewController,
		NewHandler,
	)
	return &Handler{}, nil
}
