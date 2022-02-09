package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	Id         int    `json:"id"`
	High       string `json:"high"`
	Last       string `json:"last"`
	Created_at string `json:"created_at"`
	Book       string `json:"book"`
	Volume     string `json:"volume"`
	Vwap       string `json:"vwap"`
	Low        string `json:"low"`
	Ask        string `json:"ask"`
	Bid        string `json:"bid"`
	Change_24  string `json:"change_24"`
}

type Bitso struct {
	Success bool `json:"success"`
	Payload Item `json:"payload"`
}

var url = os.Getenv("API_URL")
var dynamoTable = os.Getenv("TABLE")

func getData() {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var bitso Bitso
	json.Unmarshal([]byte(body), &bitso)
	log.Printf("%+v\n", bitso)

	if bitso.Success {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		svc := dynamodb.New(sess)
		av, err := dynamodbattribute.MarshalMap(bitso.Payload)
		if err != nil {
			log.Fatalf("Got error marshalling new item: %s", err)
		}

		tableName := dynamoTable

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			log.Fatalf("Got error calling PutItem: %s", err)
		}

		log.Printf("Successfully added to table")
	}
}

func main() {
	lambda.Start(getData)
}
