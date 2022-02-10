package services

import (
	"context"
	"bisto/internal/models"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func AddItem(db *dynamodb.Client, data models.CryptoBisto) error {	
	params := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"IdCrypto":   &types.AttributeValueMemberS{Value: data.IdCrypto},
			"created_at": &types.AttributeValueMemberS{Value: data.CreatedAt},
			"book":       &types.AttributeValueMemberS{Value: data.Book},
			"volume":     &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", data.Volume)},
			"high":       &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", data.High)},
			"last":       &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", data.Last)},
			"low":        &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", data.Low)},
			"vwap":       &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", data.Vwap)},
			"ask":        &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", data.Ask)},
			"bid":        &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", data.Bid)},
			"change_24":  &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", data.Change_24)},
			"created_on": &types.AttributeValueMemberS{Value: data.CreatedOn},
		},
		TableName: aws.String("CryptoCurrency"),
	}
	_, err := db.PutItem(context.TODO(), params)
	if err != nil {
		return err
	}
	return nil
}
