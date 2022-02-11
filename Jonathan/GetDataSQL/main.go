// main.go
//GOOS=linux GOARCH=amd64 go build -o main main.go
//zip main.zip main
package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type myInput struct {
	InitDate string `json:"initdate"`
	EndDate  string `json:"endDate"`
	Filter   string `json:"filter"`
}

type Payload struct {
	High       float64 `json:"high"`
	Last       float64 `json:"last"`
	Created_at string  `json:"created_at"`
	Book       string  `json:"book"`
	Volume     float64 `json:"volume"`
	Vwap       float64 `json:"vwap"`
	Low        float64 `json:"low"`
	Ask        float64 `json:"ask"`
	Bid        float64 `json:"bid"`
	Change_24  float64 `json:"change_24"`
}

type Item struct {
	ID      string  `json:"id"`
	Payload Payload `json:"payload"`
	Success bool    `json:"success"`
}

const (
	host     = "bisto.cp51ynp91h9l.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "arkus123"
	dbname   = "bisto"
)

func GetDataSQL(input myInput) ([]Item, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlQuery := "SELECT id, success, high, last, created_at, book, volume, vwap, low, ask, bid, change_24 FROM apibisto_data"

	rows, err := db.Query(sqlQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var Items []Item

	for rows.Next() {

		var Item_in Item

		err = rows.Scan(
			&Item_in.ID,
			&Item_in.Success,
			&Item_in.Payload.High,
			&Item_in.Payload.Last,
			&Item_in.Payload.Created_at,
			&Item_in.Payload.Book,
			&Item_in.Payload.Volume,
			&Item_in.Payload.Vwap,
			&Item_in.Payload.Low,
			&Item_in.Payload.Ask,
			&Item_in.Payload.Bid,
			&Item_in.Payload.Change_24,
		)
		if err != nil {
			panic(err)
		}
		Items = append(Items, Item_in)
	}

	// fmt.Println(Items)

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return Items, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(GetDataSQL)

}
