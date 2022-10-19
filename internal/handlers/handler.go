package handlers

import (
	"credit_holidays/internal/controllers"
)

// Handler struct for declaring api methods
type Handler struct {
	ctrl *controllers.Controller
}

func NewHandler(s *controllers.Controller) (*Handler, error) {
	return &Handler{ctrl: s}, nil
}
