package models

type Category struct {
	ID          int    `json:"category_id,omitempty"`
	Name        string `json:"category_name,omitempty"`
	Description string `json:"category_description,omitempty"`
}
