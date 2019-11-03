package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	//nvm i'm stupid

	http.HandleFunc("/repocheck/v1/commits", commitsHandler)
	//for testing
	//log.Fatal(http.ListenAndServe(":8080", nil))
	//for heroku
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
