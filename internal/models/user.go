package models

type GetBalanceRequest string

type User struct {
	Id            int64 `json:"id"`
	Balance       int64 `json:"balance"`
	FrozenBalance int64 `json:"frozen_balance"`
}
