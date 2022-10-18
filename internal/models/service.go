package models

type Service struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ConfNeeded  bool   `json:"confirmation_needed"`
	ServiceType string `json:"service_type,omitempty"`
}
