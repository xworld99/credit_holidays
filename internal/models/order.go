package models

import "time"

type Order struct {
	Id          int64
	CreatedAt   time.Time
	ProofedAt   time.Time
	UserId      int64
	ServiceId   uint
	Amount      float64
	Description string
	Status      string
}
