package models

import "time"

type AddOrderRequest struct {
	UserId    int64 `json:"user_id"`
	ServiceId int64 `json:"service_id"`
	Amount    int64 `json:"amount"`
}

type ChangeOrderRequest struct {
	OrderId int64  `json:"order_id"`
	Action  string `json:"action"`
}

type Order struct {
	Id           int64      `json:"id"`
	CreatedAt    *time.Time `json:"created_at"`
	ProofedAt    *time.Time `json:"proofed_at,omitempty"`
	ProofedAtStr string     `json:"-"`
	UserId       int64      `json:"-"`
	ServiceId    int64      `json:"-"`
	Amount       int64      `json:"amount"`
	Status       string     `json:"status"`
}
