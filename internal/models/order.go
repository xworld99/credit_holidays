package models

import "time"

type Order struct {
	Id           int64     `json:"id"`
	CreatedAt    time.Time `json:"create_at"`
	ProofedAt    time.Time `json:"proofed_at,omitempty"`
	ProofedAtStr string    `json:"-"`
	UserId       int64     `json:"-"`
	ServiceId    int64     `json:"-"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
}
