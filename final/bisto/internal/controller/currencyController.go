package controller

import (
	"bisto/internal/models"
	"bisto/internal/repository"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func notFound(c echo.Context) error {
	var response = models.ResponseCurrencies{}
	response.Success = false
	response.Message = "Currency information not found"
	jsonData, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusOK, jsonData)
}

func successResponse(c echo.Context, currencies []models.Currency, message string) error {
	if len(currencies) > 0 {
		var response = models.ResponseCurrencies{}
		response.Success = true
		response.Data = &currencies
		response.Message = message
		jsonData, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, jsonData)
	} else {
		return notFound(c)
	}
}

func validateParams(r models.RequestCurrency) []string {
	var messages []string
	const shortForm = "2006-01-02"
	if r.DateIni != "" {
		_, err := time.Parse(shortForm, r.DateIni)
		if err != nil {
			messages = append(messages, "Wrong format for dateIni, YYYY-MM-DD format is required")
		}
	}
	if r.DateEnd != "" {
		_, err := time.Parse(shortForm, r.DateEnd)
		if err != nil {
			messages = append(messages, "Wrong format for dateEnd, YYYY-MM-DD format is required")
		}
	}
	if r.Currency != "" && (r.Currency != "USD" && r.Currency != "HKD" && r.Currency != "MXN") {
		messages = append(messages, "Wrong format for Currency, only USD, HKD and MXN are accepted")
	}
	return messages
}

func GetCurrencies(c echo.Context) error {
	request := models.RequestCurrency{}
	var response = models.ResponseCurrencies{}
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	defer c.Request().Body.Close()
	//TODO: Check how to use Validate from echo and/or how check the correct params
	if err != nil {
		response.Success = false
		var errors []string
		response.Message = "Wrong parameters"
		errors = append(errors, "Provide the necessary parameters, dateInit, dateEnd in format YYYY-MM-DD and currency (USD,HKD and MXN)")
		response.Errors = errors
		jsonData, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusBadRequest, jsonData)
	}
	m := validateParams(request)
	if len(m) > 0 {
		response.Success = false
		response.Message = "Wrong parameters"
		response.Errors = m
		jsonData, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusBadRequest, jsonData)
	}
	tr := repository.NewCurrencyRepository()
	//Alls
	if request.DateIni == "" && request.DateEnd == "" && request.Currency == "" {
		currencies, _ := tr.GetAllCurrencies()
		return successResponse(c, currencies, "All currency information was converted to MXN")
	}
	//Filter by Currency
	if request.DateIni == "" && request.DateEnd == "" && request.Currency != "" {
		currencies, _ := tr.GetCurrenciesByType(request.Currency)
		msg := "All currency information was converted to " + request.Currency
		return successResponse(c, currencies, msg)
	}
	//Filter by Date
	if request.DateIni != "" && request.DateEnd != "" && request.Currency == "" {
		currencies, _ := tr.GetCurrenciesByDate(request.DateIni, request.DateEnd)
		msg := "The information of the currencies found between the dates " + request.DateIni + " and " + request.DateEnd + "was converted to MXN"
		return successResponse(c, currencies, msg)
	}
	//Filter by Date and Currency
	if request.DateIni != "" && request.DateEnd != "" && request.Currency != "" {
		currencies, _ := tr.GetCurrenciesByAllParams(request.DateIni, request.DateEnd, request.Currency)
		msg := "The information of the currencies found between the dates " + request.DateIni + " and " + request.DateEnd + "  was converted to " + request.Currency
		return successResponse(c, currencies, msg)
	}
	return notFound(c)
}
