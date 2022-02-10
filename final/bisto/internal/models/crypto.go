package models

import "time"

//github.com/webability-go/bitso - type Ticker
type CryptoBisto struct {
	IdCrypto  string  `json:"IdCrypto"`
	CreatedAt string  `json:"created_at"`
	Book      string  `json:"book"`
	Volume    float64 `json:"volume,string"`
	High      float64 `json:"high,string"`
	Last      float64 `json:"last,string"`
	Low       float64 `json:"low,string"`
	Vwap      float64 `json:"vwap,string"`
	Ask       float64 `json:"ask,string"`
	Bid       float64 `json:"bid,string"`
	Change_24 float64 `json:"change_24,string"`
	CreatedOn string  `json:"created_on"`
}

type CryptoDynamo struct {
	IdCrypto  string `json:"IdCrypto"`
	CreatedAt string `json:"created_at,string"`
	Book      string `json:"book"`
	Volume    string `json:"volume,string"`
	High      string `json:"high,string"`
	Last      string `json:"last,string"`
	Low       string `json:"low,string"`
	Vwap      string `json:"vwap,string"`
	Ask       string `json:"ask,string"`
	Bid       string `json:"bid,string"`
	Change_24 string `json:"change_24,string"`
	//CreatedOn string `json:"created_on"`
	CreatedOn *time.Time `json:"created_on"`
}
