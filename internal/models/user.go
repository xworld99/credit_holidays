package models

type User struct {
	Id            int64   `json:"id"`
	Balance       float64 `json:"balance"`
	FrozenBalance float64 `json:"frozen_balance"`
}
