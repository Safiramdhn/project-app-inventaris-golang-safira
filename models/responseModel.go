package models

type ResponseStatus bool

const (
	StatusSuccess ResponseStatus = true
	StatusError   ResponseStatus = false
)

type StandardResponse struct {
	Success ResponseStatus `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
	Errors  interface{}    `json:"errors,omitempty"`
}
