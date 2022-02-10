package models

type RequestCurrency struct {
	DateIni  string `json:"dateIni"`
	DateEnd  string `json:"dateEnd"`
	Currency string `json:"currency"`
}
