package handlers

import (
	"credit_holidays/internal/controllers"
)

// Handler struct for declaring api methods
type Handler struct {
	ctrl controllers.CreditHolidaysController
}

func NewHandler(c controllers.CreditHolidaysController) *Handler {
	return &Handler{ctrl: c}
}

func (h *Handler) Close() {
	h.ctrl.Close()
}
