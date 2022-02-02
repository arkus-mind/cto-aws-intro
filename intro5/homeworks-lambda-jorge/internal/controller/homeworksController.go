package controller

import (
	"fmt"

	"github.com/arkus-mind/intro5/homeworks-lambda-jorge/homeworksLambdaAPI/internal/service"
	"github.com/aws/aws-lambda-go/lambda"
)

func HomeworksController() {

	lambda.Start(service.GetHomeworkByTitle)
	fmt.Println("Lambda start")

}
