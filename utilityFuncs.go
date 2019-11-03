package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/api/iterator"
)

//return the authentication key
func authFinder(r *http.Request) string {
	// check if the token is in the URL
	if r.URL.Query()["private_token"] != nil {
		auth := r.URL.Query()["private_token"][0]
		return auth
	} else {
		return ""
	}
}

//returns the limit on the url
func limitFinder(r *http.Request) int {
	//Check if there's a limit
	if r.URL.Query()["limit"] != nil {
		//gets it and converts it to an integer
		customLimit, err := strconv.Atoi(r.URL.Query()["limit"][0])
		if err == nil {
			return customLimit
		} else {
			log.Fatalln(err)
		}
	} //5 is a default limit if none is present
	return 5
}

//find pages and return them
func pageFinder(auth string) int {
	page := 1
	resp, err := http.Get(fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?per_page=100&private_token=%s", auth))
	if err != nil {
		log.Fatalln(err)
	}
	//If everying is A-OK, get the pages
	if resp.StatusCode == 200 {
		page, err = strconv.Atoi(resp.Header.Get("X-Total-Pages"))
		if err != nil {
			log.Fatalln(err)
		}
	}
	return page
}

//returns page status
func statusFinder(auth string) int {
	resp, err := http.Get(fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?per_page=100&private_token=%s", auth))
	if err != nil {
		log.Fatalln(err)
	}

	return resp.StatusCode
}

// goes through all webhooks in the database and sends post
func webHooksInvoke(event string, params []string, currentTime time.Time) {
	iter := Db.Client.Collection("webhooks").Documents(Db.Ctx)
	for {
		var webHook Webhook
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate: %v", err)
		}

		err = doc.DataTo(&webHook)
		if err != nil {
			log.Printf("Failed to parse database response: %v", err)
		}

		// send POST if event is valid and is in the database
		if webHook.Event == event {
			requestWebhook := Invocation{
				Event:  event,
				Params: params,
				Time:   time.Now()}

			requestBody, err := json.Marshal(requestWebhook)
			if err != nil {
				log.Printf("Failed to marshall webhook call: %v", err)
			}

			_, err = http.Post(webHook.Url, "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				log.Printf("Not able to send webhook call: %v", err)
			}
		}
	}
}
