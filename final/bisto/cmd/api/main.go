package main

import (
	"bisto/internal/config"
	"bisto/internal/controller"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	config.InitRouting(e) //Only Local
	server := controller.WrapRouter(e)
	lambda.Start(server)	
}
