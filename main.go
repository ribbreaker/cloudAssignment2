package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hey, there's nothing on this page, try one of these instead!\n"+
		"/repocheck/v1/commits\n"+
		"/repocheck/v1/languages\n"+
		"/repocheck/v1/webhooks\n"+
		"/repocheck/v1/status", html.EscapeString(r.URL.Path))
}

func main() {

	const projectID = "cloudassignment2-f6ffb"
	const collection = "webhooks"

	Db = FirestoreDatabase{ProjectID: projectID, CollectionName: collection}
	err := Db.Init()
	if err != nil {
		log.Fatal(err)
	}

	defer Db.Close()

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/repocheck/v1/commits", CommitsHandler)
	http.HandleFunc("/repocheck/v1/languages/", LanguagesHandler)
	http.HandleFunc("/repocheck/v1/webhooks", WebHooksHandler)
	http.HandleFunc("/repocheck/v1/webhooks/", WebHooksHandler)
	http.HandleFunc("/repocheck/v1/status/", StatusHandler)

	println("listening on ", 8080)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
