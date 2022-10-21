package handlers

import (
	"credit_holidays/internal/controllers"
)

// Handler struct for declaring api methods
type Handler struct {
	ctrl *controllers.Controller
}

func NewHandler(c *controllers.Controller) *Handler {
	return &Handler{ctrl: c}
}

func (h *Handler) Close() {
	h.ctrl.Close()
}
