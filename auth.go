package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

var oauth2config oauth2.Config

func loadTokenFile() oauth2.Token {
	var token oauth2.Token

	file, err := os.Open("token.json")
	if err != nil {
		log.Fatalln("Error - unable to open token file at token.json: ", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&token)
	if err != nil {
		log.Fatalln("Error - unable to decode token JSON file: ", err)
	}

	return token
}

// Homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Homepage hit")
	u := oauth2config.AuthCodeURL(os.Getenv("stateRand"))
	http.Redirect(w, r, u, http.StatusFound)
}

// Authorize
func authorize(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := r.Form.Get("state")
	if state != os.Getenv("stateRand") {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := oauth2config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// print token to screen
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(*token)

	// print token to commandline
	fmt.Println("Token: " + token.AccessToken)

	//save token to file
	file, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		fmt.Printf("Error converting token to file: %s", err)
	}

	err = ioutil.WriteFile("token.json", file, 0644)
	if err != nil {
		fmt.Printf("Error saving token to token.json file: %s", err)
	} else {
		fmt.Println("Token successfully saved to token.json")
	}
}

func auth() {
	port := "3000"
	host := "http://localhost:" + port
	redirectEndPoint := "/oauth2"

	fmt.Println(host + redirectEndPoint)

	oauth2config = oauth2.Config{
		ClientID:     os.Getenv("clientId"),
		ClientSecret: os.Getenv("clientSecret"),
		Scopes:       []string{"read:vat", "write:vat"}, //write:vat
		// Scopes:      []string{"Read-Only"},
		RedirectURL: host + redirectEndPoint,
		// This points to our Authorization Server
		// if our Client ID and Client Secret are valid
		// it will attempt to authorize our user
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("apiUrl") + "/oauth/authorize",
			TokenURL: os.Getenv("apiUrl") + "/oauth/token",
		},
	}

	// 1 - We attempt to hit our Homepage route
	// if we attempt to hit this unauthenticated, it
	// will automatically redirect to our Auth
	// server and prompt for login credentials
	http.HandleFunc("/", homePage)

	// 2 - This displays our state, code and
	// token and expiry time that we get back
	// from our Authorization server
	http.HandleFunc(redirectEndPoint, authorize)

	// 3 - We start up our Client on port 3000
	log.Printf("Client is running at %s port.\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
