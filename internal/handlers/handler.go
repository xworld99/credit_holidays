package handlers

import (
	"credit_holidays/internal/controllers"
)

// Handler struct for declaring api methods
type Handler struct {
	ctrl controllers.CreditHolidaysController
}

// NewHandler constructor for Handler, user for code generation in wire
func NewHandler(c controllers.CreditHolidaysController) *Handler {
	return &Handler{ctrl: c}
}

// Close closes db connection
func (h *Handler) Close() {
	h.ctrl.Close()
}
