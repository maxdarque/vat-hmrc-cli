package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func readEnvFile(fileName string) {
	config := make(map[string]interface{})

	//filename is the path to the json config file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error - unable to open file: %s\n", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error - unable to decode JSON: %s\n", err)
	}

	for k, v := range config {
		// os.Setenv("uri", config["uri"])
		// os.Setenv(k, v)
		os.Setenv(k, fmt.Sprintf("%v", v))
	}
}

func main() {
	// read config file at env.json
	readEnvFile("env.json")

	//decide what to do.
	flag.Parse()
	args := flag.Args()

	// handle no args
	if len(args) == 0 {
		fmt.Println("You need to submit a command. See 'help' command")
	} else if args[0] == "help" {
		fmt.Printf("The default config file is env.json and auth/token is saved to token.json. Commands: \n\nauth - launches oauth2 authentication server on http://localhost:3000\n\nget - command used to get information from server includeing\n\t#1 'get owned' what you owe\n\t#2 'get owned-dates fromDate toDate' what you owe\n\t#3 'get liabilities fromDate toDate' for outstanding payments\n\t#4 'get payments fromDate toDate' for payments HMRC has received\n\t#5 'get return periodKey' get the return for a given period denoted by the periodKey i.e. 18AD\n[Note: dates to be in following format: yyyy-mm-dd. If no date is given, then the program will default to the preciding one year period. The 'get' command can use numbers shown instead of commands i.e. 'get 2']\n\nsubmit - 'submit fileName.json' see example json file\n")
	} else if args[0] == "auth" {
		auth()
	} else if args[0] == "get" {
		get(args)
	} else if args[0] == "submit" {
		submit(args)
	} else {
		fmt.Println("Command not recognized")
	}
}
