package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

//vatReturn - Structure for a VAT return
type vatReturn struct {
	PeriodKey                    string  `json:"periodKey"`
	VatDueSales                  float64 `json:"vatDueSales"`
	VatDueAcquisitions           float64 `json:"vatDueAcquisitions"`
	TotalVatDue                  float64 `json:"totalVatDue"`
	VatReclaimedCurrPeriod       float64 `json:"vatReclaimedCurrPeriod"`
	NetVatDue                    float64 `json:"netVatDue"`
	TotalValueSalesExVAT         float64 `json:"totalValueSalesExVAT"`
	TotalValuePurchasesExVAT     float64 `json:"totalValuePurchasesExVAT"`
	TotalValueGoodsSuppliedExVAT float64 `json:"totalValueGoodsSuppliedExVAT"`
	TotalAcquisitionsExVAT       float64 `json:"totalAcquisitionsExVAT"`
	Finalised                    bool    `json:"finalised"`
}

// {
//   "periodKey": "A001",
//   "vatDueSales": 105.50,
//   "vatDueAcquisitions": -100.45,
//   "totalVatDue": 5.05,
//   "vatReclaimedCurrPeriod": 105.15,
//   "netVatDue": 100.10,
//   "totalValueSalesExVAT": 300,
//   "totalValuePurchasesExVAT": 300,
//   "totalValueGoodsSuppliedExVAT": 3000,
//   "totalAcquisitionsExVAT": 3000,
//   "finalised": true
// }

func submit(args []string) {
	if len(args) != 2 {
		log.Fatalln("Arguments required 'submit filename.json")
	}

	fileName := args[1]
	var v vatReturn

	//load json
	//fileName is the path to the json config file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error - unable to open file: %s\n", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&v)
	if err != nil {
		log.Fatalf("Error - unable to decode JSON: %s\n", err)
	}

	//convert struct to []byte
	vatReturnJSON, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Error - unable to convert JSON to byte array: %s\n", err)
	}

	// fmt.Printf("K %f\n", vatReturnJSON.NetVatDue)

	//create the url and send the post with the details
	u := os.Getenv("apiUrl") + "/organisations/vat/" + os.Getenv("vrn") + "/returns"

	fmt.Printf("Post URL: " + u + "\n\n")

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)

	client := &http.Client{}
	request, _ := http.NewRequest("POST", u, bytes.NewBuffer(vatReturnJSON))
	setHeaders(request)

	response, err := client.Do(request)

	if err != nil {
		log.Fatalf("HTTP post request failed with code: %s\n", err)
	}
	defer response.Body.Close()

	saveResponse(response)
}
