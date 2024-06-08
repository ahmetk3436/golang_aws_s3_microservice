package domain

type Response struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}
