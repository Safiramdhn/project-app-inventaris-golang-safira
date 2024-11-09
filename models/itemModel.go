package models

import "time"

type Item struct {
	ID                  int       `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	CategoryID          int       `json:"category_id,omitempty"`
	CategoryName        string    `json:"category,omitempty"`
	PhotoURL            string    `json:"photo_url,omitempty"`
	Price               float64   `json:"price,omitempty"`
	PurchaseDate        time.Time `json:"purchase_date,omitempty"`
	TotalUsageDays      int       `json:"total_usage_days,omitempty"`
	IsReplacementNeeded bool      `json:"is_replacement_needed,omitempty"`
	DepreciatedRate     int       `json:"depresiated_rate,omitempty"`
}
