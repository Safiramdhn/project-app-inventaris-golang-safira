package models

import "time"

type ItemInvestment struct {
	ID                   int       `json:"id,omitempty"`
	ItemID               int       `json:"item_id,omitempty"`
	ItemName             string    `json:"item_name,omitempty"`
	InitialPrice         float64   `json:"initial_price,omitempty"`
	CurrentValue         float64   `json:"current_value,omitempty"`
	DepreciationRate     float64   `json:"depreciation_rate,omitempty"`
	LastDepreciationDate time.Time `json:"last_depreciation_date,omitempty"`
	TotalInvestment      float64   `json:"total_investment,omitempty"`
	DepricatedValue      float64   `json:"depricated_value,omitempty"`
}
