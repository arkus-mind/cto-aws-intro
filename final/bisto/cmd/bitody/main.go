package main

import (
	"bisto/internal/controller"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(controller.WriteToDynamo)
}
