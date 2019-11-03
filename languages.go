package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func LanguagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		auth := authFinder(r)
		limit := limitFinder(r)
		page := pageFinder(auth)

		var idStruct []IdStruct
		ranking := make(map[string]int)
		var langStruct Languages

		// Payload handling
		var payload []string
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err == io.EOF {
			payload = nil
		} else if err != nil {
			http.Error(w, "Payload is bad", http.StatusBadRequest)
			return
		}

		// invoking the webhook for those watching language
		params := []string{strconv.Itoa(limit)}
		webHooksInvoke("languages", params, time.Now())

		//if we entered the matrix
		if statusFinder(auth) == 200 {

			for i := 1; i <= page; i++ {

				// loops through all the pages that are found
				resp, err := http.Get(fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?per_page=100&page=%d&private_token=%s", i, auth))
				if err != nil {
					log.Fatalln(err)
				}

				// decodes into the struct
				err = json.NewDecoder(resp.Body).Decode(&idStruct)
				if err != nil {
					log.Fatalln(err)
				}

				// goes through all the pages that are found
				for j := 0; j < len(idStruct); j++ {
					// prepares payload check
					isInPayload := false
					if payload != nil {
						for _, payloadItem := range payload { // goes through all things in the payload
							if idStruct[j].Name == payloadItem { // if the name in the struct and the payload match
								isInPayload = true // set flag as true
								break              // and break out of the area
							}
						}
					}

					// If there is no payload or repository was in payload then get languages, handle errors as needed.
					if payload == nil || isInPayload {
						// sets id to id from page gotten
						id := idStruct[j].Id
						resp, err = http.Get(fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%d/languages?private_token=%s", id, auth))
						if err != nil {
							log.Fatalln(err)
						}

						languages := make(map[string]interface{})

						// decodes languages into the map made above
						err = json.NewDecoder(resp.Body).Decode(&languages)
						if err != nil {
							http.Error(w, "Error decoding GitLab languages response.", http.StatusInternalServerError)
							return
						}

						for key := range languages {
							//does the actual counting of code names
							ranking[key]++

						}
					}
				}
			}
			//iterates through the limit
			for i := 0; i < limit; i++ {

				bigKey := ""
				bigValue := 0

				// loops through the map and checks values
				for key, value := range ranking {
					if value > bigValue {
						bigValue = value
						bigKey = key // pretty simple, finds the biggest value
					}
				}
				// puts the biggest value in the language array struct
				langStruct.Languages = append(langStruct.Languages, bigKey)
				delete(ranking, bigKey)
			}

			// sets auth bool to proper value
			if auth == "" {
				langStruct.Auth = false
			} else {
				langStruct.Auth = true
			}

			//adds header and encodes the languages
			w.Header().Add("content-type", "application/json")
			err := json.NewEncoder(w).Encode(langStruct)
			if err != nil {
				log.Fatalln(err)
			}

		} else {
			fmt.Fprintf(w, "Malformed request")
		}
	} else {
		println("if get - else start")
		fmt.Fprintf(w, "Get is the only method that works")
	}
}
