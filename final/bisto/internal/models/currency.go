package models

import "time"

type Currency struct {
	Id        string    `json:"Id"`
	IdCrypto  string    `json:"IdCrypto"`
	CreatedAt time.Time `json:"created_at"`
	Book      string    `json:"book"`
	Volume    float64   `json:"volume,string"`
	High      float64   `json:"high,string"`
	Last      float64   `json:"last,string"`
	Low       float64   `json:"low,string"`
	Vwap      float64   `json:"vwap,string"`
	Ask       float64   `json:"ask,string"`
	Bid       float64   `json:"bid,string"`
	Change_24 float64   `json:"change_24,string"`
	USDToMXN  float64   `json:"usd_mxn,string"`
	HKDToMXN  float64   `json:"hkd_mxn,string"`
}
