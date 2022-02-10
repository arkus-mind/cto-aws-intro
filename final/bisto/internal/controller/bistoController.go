package controller

import (
	"bisto/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const bistoURL = "https://api.bitso.com/v3/ticker"
const usd_mxn = "usd_mxn"
const btc_mxn = "btc_mxn"
const HKDToMXN = 2.64
const MXNToHKD = 0.38

func GetUSDToMXN() models.CryptoBisto {
	var cripto = models.CryptoBisto{}
	response, err := http.Get(bistoURL + "/?book=" + usd_mxn)

	if err != nil {
		fmt.Print(err.Error())
		return cripto
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject models.Response
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("error Unmarshal: ", e)
	}
	return responseObject.Payload
}

func GetDataFromBisto() models.CryptoBisto {
	var cripto = models.CryptoBisto{}
	response, err := http.Get(bistoURL + "/?book=" + btc_mxn)

	if err != nil {
		fmt.Print(err.Error())
		return cripto
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
	var responseObject models.Response
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("error Unmarshal: ", e)
	}
	return responseObject.Payload
}
