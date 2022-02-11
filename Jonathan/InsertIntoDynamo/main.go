package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/lambda"

	"log"
)

type Payload struct {
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

type Item struct {
	ID      string  `json:"id"`
	Payload Payload `json:"payload"`
	Success bool    `json:"success"`
}

func main() {
	lambda.Start(GetData)
}

func GetData() {

	api := "https://api.bitso.com/v3/ticker/?book=btc_mxn"

	//get resquest
	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//read resquest
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//create new item
	item := Item{}

	//Fill new item
	err = json.Unmarshal(body, &item)

	if err != nil {
		log.Fatalf("Cannot get de data from bisto: %s", err)
	}

	//add ID to item
	u, err := uuid.NewUUID()

	if err != nil {
		log.Fatalf("Got error creating UUID: %s", err)
	}

	id := u.String()

	item.ID = id

	//Insert
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	tableName := "apibisto_data"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	} else {
		fmt.Printf("UUID create: %s", item.ID)
	}
}
