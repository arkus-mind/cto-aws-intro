package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

func handleRequest(e events.DynamoDBEvent) {
	uri := os.Getenv("URI")

	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, record := range e.Records {
		/*
			for name, value := range record.Change.NewImage {
				log.Printf("%s", name)
				if value.DataType() == events.DataTypeString {
					fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())
				}
			}
			log.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)
			if record.Change.NewImage["Created_at"].DataType() == events.DataTypeNumber {
				log.Printf("asdf")
			}
			if record.Change.NewImage["Created_at"].DataType() == events.DataTypeString {
				log.Printf("Simona la cacarisaStr")
			} */

		date := record.Change.NewImage["created_at"].String()
		mxn, _ := strconv.ParseFloat(record.Change.NewImage["last"].String(), 64)
		hgd := mxn * 0.38
		usd := mxn * 0.048

		query := `
		INSERT INTO public.bitcoin(
			mxn, usd, hkd, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id`

		var id string = ""
		err = db.QueryRow(query, mxn, hgd, usd, date).Scan(&id)
		if err != nil {
			panic(err)
		}
		log.Printf("Record inserted with id: %s", id)
	}
}
func main() {
	lambda.Start(handleRequest)
}
