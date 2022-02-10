package config

import (
	"bisto/internal/controller"

	"github.com/labstack/echo"
)

func InitRouting(e *echo.Echo) {
	e.POST("/currency/getCurrencies", controller.GetCurrencies)
}
