package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func setHeaders(request *http.Request) {
	token := loadTokenFile()
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Set("Accept", "application/vnd.hmrc.1.0+json")
	request.Header.Set("Content-Type", "application/json")

	//Fraud prevention Headers
	deviceID := loadOrCreateDeviceID()
	//https://developer.service.hmrc.gov.uk/api-documentation/docs/fraud-prevention
	request.Header.Set("Gov-Client-Connection-Method", "DESKTOP_APP_DIRECT")
	request.Header.Set("Gov-Client-Device-ID", deviceID)
	// request.Header.Set("Gov-Client-User-IDs") not required as it can be blank or omitted according to the docs
	request.Header.Set("Gov-Client-Timezone", "UTC+00:00")
	// request.Header.Set("Gov-Client-Local-IPs","") not required as it can be blank or omitted according to the docs
	// request.Header.Set("Gov-Client-MAC-Addresses","") not required as it can be blank or omitted according to the docs
	request.Header.Set("Gov-Client-Screens", "width=1920&height=1080") //doesn't include scaling-factor or colour-depth which can be omitted
	// request.Header.Set("Gov-Client-Window-Size","") If the call was initiated by a background process, submit the header with an empty value or omit it entirely.
	request.Header.Set("Gov-Client-User-Agent", "")   //not sure if I can leave this blank
	request.Header.Set("Gov-Client-Multi-Factor", "") //If only a single factor (for example, username and password) is being used, submit the header with an empty value or omit it entirely.
	request.Header.Set("Gov-Vendor-Version", "0.01")
	request.Header.Set("Gov-Vendor-License-IDs", "") //If there are no such licenses on the originating device, then submit the header with an empty value or omit it entirely.

	//TESTING purposes
	// request.Header.Set("Gov-Test-Scenario", "SINGLE_PAYMENT_2018_19")
}

type deviceIDJSON struct {
	ID string `json:"id"`
}

func loadOrCreateDeviceID() string {
	var idJSON deviceIDJSON

	file, err := os.Open("deviceID.json")
	if err != nil {
		//create ID
		uuidObj, _ := uuid.NewRandom()
		id := uuidObj.String()

		idJSON = deviceIDJSON{ID: id}
		//save id to file
		file, err := json.MarshalIndent(idJSON, "", "  ")
		if err != nil {
			fmt.Printf("Error converting ID to file: %s", err)
		}

		err = ioutil.WriteFile("deviceID.json", file, 0644)
		if err != nil {
			fmt.Printf("Error saving deviceID to deviceID.json file: %s", err)
		} else {
			fmt.Println("DeviceID successfully saved to deviceID.json")
		}
		//return id
		return id
	}
	//load exisiting ID
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&idJSON)
	if err != nil {
		log.Fatalln("Error - unable to decode device ID JSON file: ", err)
	}
	return idJSON.ID
}

func saveResponse(response *http.Response) {
	fmt.Println("Response Status:", response.Status)
	// fmt.Println("Response Headers:", response.Header)

	//convert the response to a JSON in order to save and print to screen.
	bodyJSON := make(map[string]interface{})
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&bodyJSON)
	if err != nil {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Raw JSON: " + string(data) + "\n")
		log.Fatalf("Error - unable to decode JSON: %s\n", err)
	}

	//print JSON to the screen
	body, err := json.MarshalIndent(bodyJSON, "", "  ")
	if err != nil {
		log.Fatalf("formatting JSON: %s\n", err)
	}

	fmt.Printf("Response Body: %s\n\n", body)

	//Save response to file
	//does the logs folder exist? If not, create one.
	dir := "./logs"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	t := time.Now().Local()
	fileName := t.String() + " - Response.json"

	err = ioutil.WriteFile(dir+"/"+fileName, body, 0644)
	if err != nil {
		fmt.Printf("Error saving response to file: %s", err)
	}
	fmt.Printf("\n\nResponse saved to 'logs' folder\n")
}
