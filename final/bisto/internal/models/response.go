package models

type ResponseCurrencies struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *[]Currency `json:"data"`
	Errors  []string    `json:"errors"`
}

type Response struct {
	Success bool        `json:"success"`
	Payload CryptoBisto `json:"payload"`
}
