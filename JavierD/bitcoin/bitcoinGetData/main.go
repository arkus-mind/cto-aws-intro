package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

const limit string = "10"

type Bitcoin struct {
	ID        string  `json:"id"`
	MXN       float64 `json:"mxn"`
	USD       float64 `json:"usd"`
	HKD       float64 `json:"hkd"`
	CreadetAt string  `json:"created_at"`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {
	case "GET":
		return get_bitcoin(req)
	default:
		apiResp := events.APIGatewayProxyResponse{Body: req.HTTPMethod, StatusCode: 200}
		return apiResp, nil
	}
}

func get_bitcoin(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	initialDate := req.QueryStringParameters["iDate"]
	finalDate := req.QueryStringParameters["fDate"]

	offset := req.QueryStringParameters["offset"]
	if offset == "" {
		offset = "0"
	}

	query := "SELECT id, mxn, usd, hkd, created_at FROM public.bitcoin WHERE created_at BETWEEN $1::timestamp AND $2::timestamp ORDER BY id LIMIT $3 OFFSET $4 "

	uri := os.Getenv("URI")
	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(query, initialDate, finalDate, limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var bitcoins []Bitcoin

	for rows.Next() {
		var sBit Bitcoin
		err := rows.Scan(&sBit.ID, &sBit.MXN, &sBit.USD, &sBit.HKD, &sBit.CreadetAt)
		if err != nil {
			log.Fatal(err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			}, nil
		}
		bitcoins = append(bitcoins, sBit)
	}

	js, err := json.Marshal(bitcoins)
	log.Printf(string(js))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func main() {
	lambda.Start(router)
}
