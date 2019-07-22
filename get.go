package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func getRequest(u string) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", u, nil)
	setHeaders(request)
	fmt.Printf("Get URL: " + u + "\n\n")
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("HTTP failed with code: %s\n", err)
	}
	defer response.Body.Close()

	saveResponse(response)
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
	rootURL := os.Getenv("API_URL") + "/organisations/vat/" + os.Getenv("VRN")
	command := args[1]

	if command == "help" {
		fmt.Printf("#1 'get owned' what you owe\n#2 'get owned-dates fromDate toDate' what you owe\n#3 'get liabilities fromDate toDate' for outstanding payments\n#4 'get payments fromDate toDate' for payments HMRC has received\n#5 'get return periodKey' get the return for a given period denoted by the periodKey i.e. 18AD\n\n[Note: dates to be in following format: yyyy-mm-dd. If no date is given, then the program will default to the preciding one year period. The 'get' command can use numbers shown instead of commands i.e. 'get 2']\n")
	} else if command == "owned" || command == "1" {
		fmt.Println("Getting open VAT obligations ...")
		u = rootURL + "/obligations"
		query = "?status=O"
		getRequest(u + query)
	} else if command == "owned-dates" || command == "2" {
		fmt.Println("Getting VAT obligations within given dates ...")
		u = rootURL + "/obligations"
		query = getQueryOnDates(args)
		getRequest(u + query)
	} else if command == "liabilities" || command == "3" {
		fmt.Println("Getting outstanding payments due to HMRC...")
		u = rootURL + "/liabilities"
		query = getQueryOnDates(args)
		getRequest(u + query)
	} else if command == "payments" || command == "4" {
		fmt.Println("Getting payments HMRC has received...")
		u = rootURL + "/payments"
		query = getQueryOnDates(args)
		getRequest(u + query)
	} else if command == "return" || command == "5" {
		fmt.Println("Getting return for given period...")
		if len(args) != 3 {
			log.Fatalln("Period key needed to fetch return. See help")
		}
		periodKey := url.QueryEscape(args[2])
		u = rootURL + "/returns/" + periodKey
		getRequest(u)
	} else {
		fmt.Println("Command not recognized")
	}
}
