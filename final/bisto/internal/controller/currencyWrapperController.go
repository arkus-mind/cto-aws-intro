package controller

import (
	"bisto/internal/models"
	"bisto/internal/repository"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo"
)

// WrapEchoServer wraps echo server into Lambda Handler
func WrapRouter(e *echo.Echo) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var response = models.ResponseCurrencies{}
		body := strings.NewReader(request.Body)
		dataRequest := models.RequestCurrency{}
		err := json.NewDecoder(body).Decode(&dataRequest)
		//TODO: Check how to use Validate from APIGatewayProxyRequest and/or how check the correct params
		if err != nil {
			response.Success = false
			var errors []string
			response.Message = "Wrong parameters"
			errors = append(errors, "Provide the necessary parameters, dateInit, dateEnd in format YYYY-MM-DD and currency (USD,HKD and MXN)")
			response.Errors = errors
			jsonData, _ := json.Marshal(response)
			return formatAPIErrorResponse(http.StatusBadRequest, httptest.NewRecorder().Header(), string(jsonData))
		}

		m := validateParams(dataRequest)

		if len(m) > 0 {
			response.Success = false
			response.Message = "Wrong parameters"
			response.Errors = m
			jsonData, _ := json.Marshal(response)
			return formatAPIErrorResponse(http.StatusBadRequest, httptest.NewRecorder().Header(), string(jsonData))
		}
		tr := repository.NewCurrencyRepository()
		//Alls
		if dataRequest.DateIni == "" && dataRequest.DateEnd == "" && dataRequest.Currency == "" {
			currencies, _ := tr.GetAllCurrencies()
			return formatAPIResponse(http.StatusOK, httptest.NewRecorder().Header(), currencies, "All currency information was converted to MXN")
		}
		//Filter by Currency
		if dataRequest.DateIni == "" && dataRequest.DateEnd == "" && dataRequest.Currency != "" {
			currencies, _ := tr.GetCurrenciesByType(dataRequest.Currency)
			msg := "All currency information was converted to " + dataRequest.Currency
			return formatAPIResponse(http.StatusOK, httptest.NewRecorder().Header(), currencies, msg)
		}
		//Filter by Date
		if dataRequest.DateIni != "" && dataRequest.DateEnd != "" && dataRequest.Currency == "" {
			currencies, _ := tr.GetCurrenciesByDate(dataRequest.DateIni, dataRequest.DateEnd)
			msg := "The information of the currencies found between the dates " + dataRequest.DateIni + " and " + dataRequest.DateEnd + "was converted to MXN"
			return formatAPIResponse(http.StatusOK, httptest.NewRecorder().Header(), currencies, msg)
		}
		//Filter by Date and Currency
		if dataRequest.DateIni != "" && dataRequest.DateEnd != "" && dataRequest.Currency != "" {
			currencies, _ := tr.GetCurrenciesByAllParams(dataRequest.DateIni, dataRequest.DateEnd, dataRequest.Currency)
			msg := "The information of the currencies found between the dates " + dataRequest.DateIni + " and " + dataRequest.DateEnd + "  was converted to " + dataRequest.Currency
			return formatAPIResponse(http.StatusOK, httptest.NewRecorder().Header(), currencies, msg)
		}

		return formatAPIErrorResponse(http.StatusBadRequest, httptest.NewRecorder().Header(), "Currency information not found")
	}
}

func formatAPIResponse(statusCode int, headers http.Header, currencies []models.Currency, message string) (events.APIGatewayProxyResponse, error) {
	if len(currencies) > 0 {
		responseHeaders := make(map[string]string)

		responseHeaders["Content-Type"] = "application/json"
		for key, value := range headers {
			responseHeaders[key] = ""

			if len(value) > 0 {
				responseHeaders[key] = value[0]
			}
		}

		responseHeaders["Access-Control-Allow-Origin"] = "*"
		responseHeaders["Access-Control-Allow-Headers"] = "origin,Accept,Authorization,Content-Type"
		var response = models.ResponseCurrencies{}
		response.Success = true
		response.Data = &currencies
		response.Message = message
		jsonData, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			Body:       string(jsonData),
			Headers:    responseHeaders,
			StatusCode: statusCode,
		}, nil
	}
	return formatAPIErrorResponse(http.StatusBadRequest, headers, "Currency information not found")
}

func formatAPIErrorResponse(statusCode int, headers http.Header, err string) (events.APIGatewayProxyResponse, error) {
	responseHeaders := make(map[string]string)

	responseHeaders["Content-Type"] = "application/json"
	for key, value := range headers {
		responseHeaders[key] = ""

		if len(value) > 0 {
			responseHeaders[key] = value[0]
		}
	}

	responseHeaders["Access-Control-Allow-Origin"] = "*"
	responseHeaders["Access-Control-Allow-Headers"] = "origin,Accept,Authorization,Content-Type"
	var response = models.ResponseCurrencies{}
	response.Success = false
	response.Message = err
	jsonData, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		Body:       string(jsonData),
		Headers:    responseHeaders,
		StatusCode: statusCode,
	}, nil
}
