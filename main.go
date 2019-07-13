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
	fmt.Println()
	// read config file at env.json
	readEnvFile("env.json")

	//decide what to do.
	flag.Parse()
	args := flag.Args()

	// handle no args
	if len(args) == 0 {
		fmt.Println("You need to submit a command. See 'help' command")
	} else if args[0] == "help" {
		fmt.Printf("The default config file is env.json and auth/token is saved to token.json. Commands: \nauth - launches oauth2 authentication server on http://localhost:3000\nget - command used to get information from server. type 'get help' for more info.\nsubmit - 'submit fileName.json' see example json file\n\n")
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
