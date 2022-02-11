// main.go
//GOOS=linux GOARCH=amd64 go build -o main main.go
//zip main.zip main
package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"

	_ "github.com/lib/pq"
)

const (
	host     = "bisto.cp51ynp91h9l.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "arkus123"
	dbname   = "bisto"
)

type Item struct {
	ID      string  `json:"id"`
	Payload Payload `json:"payload"`
	Success bool    `json:"success"`
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

func DynamoToPosgresql(ctx context.Context, e events.DynamoDBEvent) (string, error) {

	item := Item{}

	for _, record := range e.Records {

		item.ID = record.Change.NewImage["id"].String()

		payload := record.Change.NewImage["payload"].Map()

		item.Payload.High = convertStringToFloat64(payload["high"].String())
		item.Payload.Last = convertStringToFloat64(payload["last"].String())
		item.Payload.Created_at = payload["created_at"].String()
		item.Payload.Book = payload["book"].String()
		item.Payload.Volume = convertStringToFloat64(payload["volume"].String())
		item.Payload.Vwap = convertStringToFloat64(payload["vwap"].String())
		item.Payload.Low = convertStringToFloat64(payload["low"].String())
		item.Payload.Ask = convertStringToFloat64(payload["ask"].String())
		item.Payload.Bid = convertStringToFloat64(payload["bid"].String())
		item.Payload.Change_24 = convertStringToFloat64(payload["change_24"].String())

		item.Success = record.Change.NewImage["success"].Boolean()

		fmt.Print(item)

	}

	return insertIntoDB(item), nil
}

func insertIntoDB(item Item) string {
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

	sqlStatement := `
	INSERT INTO public.apibisto_data (
		id, success, high, last, created_at, book, volume, vwap, low, ask, bid, change_24
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	)
	RETURNING id`
	id := ""
	err = db.QueryRow(sqlStatement,
		item.ID,
		item.Success,
		item.Payload.High,
		item.Payload.Last,
		item.Payload.Created_at,
		item.Payload.Book,
		item.Payload.Volume,
		item.Payload.Vwap,
		item.Payload.Low,
		item.Payload.Ask,
		item.Payload.Bid,
		item.Payload.Change_24).Scan(&id)
	if err != nil {
		panic(err)
	}
	return "New record ID is: " + id
}

func convertStringToFloat64(intro string) float64 {

	result, err := strconv.ParseFloat(intro, 64)
	if err != nil {
		panic(err)
	}
	return result
}

func main() {
	lambda.Start(DynamoToPosgresql)
}
