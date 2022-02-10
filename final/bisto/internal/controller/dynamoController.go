package controller

import (	
	"bisto/internal/repository"
	"bisto/internal/services"
	"bisto/internal/models"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func WriteToDynamo() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	data := GetDataFromBisto()
	if data.Book != "" {
		data.IdCrypto = repository.GenerateUUID().String()
		data.CreatedOn = time.Now().Format(time.RFC3339)
		err := services.AddItem(svc, data)
		if err != nil {
			fmt.Println("failed Add Item to DynamoDB, ", err)
		}
	}
}

func ReadFromDynamoWriteToRDS() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)
	now := time.Now()
	fmt.Println("Now: ", now)
	fmt.Println()
	//TODO: Check what best expression I can use to avoid losing information
	filt := expression.Name("created_on").LessThanEqual(expression.Value(now.Format(time.RFC3339)))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("CryptoCurrency"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(context.TODO(), params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}

	numItems := 0
	currentChange := GetUSDToMXN()

	cr := repository.NewCurrencyRepository()
	fmt.Println("Found Items: ", len(result.Items))
	for _, i := range result.Items {
		item := models.CryptoDynamo{}
		currency := models.Currency{}

		err = attributevalue.UnmarshalMap(i, &item)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}

		if item.IdCrypto != "" {
			if cr.ExistCurrency(item.IdCrypto) {
				continue
			}
			numItems++
			currency.Id = repository.GenerateUUID().String()
			currency.IdCrypto = item.IdCrypto
			currency.Book = item.Book
			currency.CreatedAt = time.Now()
			currency.Volume, _ = strconv.ParseFloat(item.Volume, 64)
			currency.High, _ = strconv.ParseFloat(item.High, 64)
			currency.Last, _ = strconv.ParseFloat(item.Last, 64)
			currency.Low, _ = strconv.ParseFloat(item.Low, 64)
			currency.Vwap, _ = strconv.ParseFloat(item.Vwap, 64)
			currency.Ask, _ = strconv.ParseFloat(item.Ask, 64)
			currency.Bid, _ = strconv.ParseFloat(item.Bid, 64)
			currency.Change_24, _ = strconv.ParseFloat(item.Change_24, 64)
			currency.USDToMXN = currentChange.Ask
			currency.HKDToMXN = HKDToMXN
			fmt.Println()
			id := cr.NewCurrency(currency)
			fmt.Println("Id Currency: ", id)
			fmt.Println()
		}
	}
	fmt.Println("Found new items: ", numItems)
	cr.CloseConnection()
}
