package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CommitsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		auth := authFinder(r)
		limit := limitFinder(r)
		page := pageFinder(auth)

		var commit Commit
		var repo Repository
		var repoId []IdStruct
		commitsNum := make(map[int]int)

		// invokes the webhook for all those that are listening on commits
		params := []string{strconv.Itoa(limit)}
		webHooksInvoke("commits", params, time.Now())

		if statusFinder(auth) == 200 {
			// goes through the number of repoes specified in limit
			for i := 1; i <= page; i++ {

				// loops through all the pages that are found
				resp, err := http.Get(fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?per_page=100&page=%d&private_token=%s", i, auth))
				if err != nil {
					log.Fatalln(err)
				}

				// decodes into the struct
				err = json.NewDecoder(resp.Body).Decode(&repoId)
				if err != nil {
					log.Fatalln(err)
				}

				for j := 0; j < len(repoId); j++ {

					// sets id to id from page gotten
					id := repoId[j].Id
					resp, err = http.Get(fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%d/repository/commits?private_token=%s", id, auth))
					if err != nil {
						log.Fatalln(err)
					}

					commits, err := strconv.Atoi(resp.Header.Get("X-Total"))
					if err != nil {
						println("error in ", id)
					}
					commitsNum[id] = commits
				}
			}

			for i := 0; i < limit; i++ {

				// Makes variables
				bigKey := 0
				bigValue := 0

				// loops through the map and checks values
				for key, value := range commitsNum {
					if value > bigValue {
						bigValue = value
						bigKey = key // pretty simple, finds the biggest value
					}
				}
				resp, err := http.Get(fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%d?private_token=%s", bigKey, auth))
				if err != nil {
					log.Fatalln(err)
				}

				err = json.NewDecoder(resp.Body).Decode(&repo)
				if err != nil {
					log.Fatalln(err)
				}
				repo.Commits = bigValue

				// puts the biggest value in the language array struct
				commit.Repos = append(commit.Repos, repo)
				delete(commitsNum, bigKey)
			}

			// sets auth bool to proper value
			if auth == "" {
				commit.Auth = false
			} else {
				commit.Auth = true
			}

			//adds header and encodes the languages
			w.Header().Add("content-type", "application/json")
			err := json.NewEncoder(w).Encode(commit)
			if err != nil {
				log.Fatalln(err)
			}

		} else {
			fmt.Fprintf(w, "Malformed request")
		}
	}
}
