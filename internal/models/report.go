package models

import (
	"fmt"
	"time"
)

type GetHistoryRequest struct {
	UserId   string
	FromDate string
	ToDate   string
	Offset   string
	Limit    string
	OrderBy  string
}

type History struct {
	OrderId            int64      `json:"order_id"`
	CreatedAt          *time.Time `json:"created_at"`
	ProofedAt          *time.Time `json:"proofed_at,omitempty"`
	Status             string     `json:"status"`
	ServiceName        string     `json:"service_name"`
	ServiceDescription string     `json:"service_description"`
	ServiceType        string     `json:"service_type"`
	Amount             int64      `json:"amount"`
}

type HistoryFrame struct {
	TotalOperations int       `json:"total_operations"`
	Offset          int       `json:"current_offset"`
	Limit           int       `json:"-"`
	FromDate        time.Time `json:"-"`
	ToDate          time.Time `json:"-"`
	OrderBy         string    `json:"-"`
	UserId          int64     `json:"user_id"`
	Operations      []History `json:"operations"`
}

type GenerateReportRequest string

type CSVRow struct {
	Id       int64
	Name     string
	Type     string
	CashFlow int64
}

type CSVData struct {
	Period  time.Time
	Records []CSVRow
}

func (c *CSVData) ToStringSlice() [][]string {
	res := [][]string{{"ServiceId", "Name", "Type", "CashFlow"}}

	for _, r := range c.Records {
		res = append(res, []string{
			fmt.Sprintf("%d", r.Id),
			r.Name,
			r.Type,
			fmt.Sprintf("%d", r.CashFlow),
		})
	}

	return res
}
