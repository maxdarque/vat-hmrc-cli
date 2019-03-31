package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func getRequest(u string) {
	token := loadTokenFile()
	client := &http.Client{}
	request, _ := http.NewRequest("GET", u, nil)
	request.Header.Set("Accept", "application/vnd.hmrc.1.0+json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	fmt.Printf("Get URL: " + u + "\n\n")
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("HTTP failed with code: %s\n", err)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	//convert the response to a pretty JSON to print to screen.
	bodyJson := make(map[string]interface{})
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&bodyJson)
	if err != nil {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Raw JSON: " + string(data) + "\n")
		log.Fatalf("Error - unable to decode JSON: %s\n", err)
	}

	body, err := json.MarshalIndent(bodyJson, "", "  ")
	if err != nil {
		log.Fatalf("Error converting tokent to file: %s\n", err)
	}

	fmt.Printf("Response body: %s\n\n", body)

	// data, _ := ioutil.ReadAll(response.Body)
	// fmt.Println(string(data))
}

func getQueryOnDates(args []string) string {
	var from, to string
	if len(args) >= 4 {
		from = args[2]
		to = args[3]
	} else {
		now := time.Now()
		to = now.Format("2006-01-02")
		from = now.AddDate(-1, 0, 1).Format("2006-01-02")
		fmt.Printf("Warning: This program will request data for the last year (%s to %s) unless otherwise specified. Use the third and fourth argument for the from and to dates in format YYYY-MM-DD\n\n", from, to)
	}
	return "?from=" + from + "&to=" + to
}

func get(args []string) {

	if len(args) < 2 {
		log.Fatalln("More commands needed. See 'help' command")
	}

	var u, query string
	rootUrl := os.Getenv("apiUrl") + "/organisations/vat/" + os.Getenv("vrn")
	command := args[1]

	if command == "owned" || command == "1" {
		fmt.Println("Getting open VAT obligations ...")
		u = rootUrl + "/obligations"
		query = "?status=O"
		getRequest(u + query)
	} else if command == "owned-dates" || command == "2" {
		fmt.Println("Getting VAT obligations within given dates ...")
		u = rootUrl + "/obligations"
		query = getQueryOnDates(args)
		getRequest(u + query)
	} else if command == "liabilities" || command == "3" {
		fmt.Println("Getting outstanding payments due to HMRC...")
		u = rootUrl + "/liabilities"
		query = getQueryOnDates(args)
		getRequest(u + query)
	} else if command == "payments" || command == "4" {
		fmt.Println("Getting payments HMRC has received...")
		u = rootUrl + "/payments"
		query = getQueryOnDates(args)
		getRequest(u + query)
	} else if command == "return" || command == "5" {
		fmt.Println("Getting return for given period...")
		if len(args) != 3 {
			log.Fatalln("Period key needed to fetch return. See help")
		}
		periodKey := url.QueryEscape(args[2])
		u = rootUrl + "/returns/" + periodKey
		getRequest(u)
	} else {
		fmt.Println("Command not recognized")
	}
}
