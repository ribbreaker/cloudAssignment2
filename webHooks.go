package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var Db FirestoreDatabase

func WebHooksHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

		//GET
	case http.MethodGet:

		webHookId := ""
		webHookId = r.URL.Path[23:]
		var iteration *firestore.DocumentIterator
		var webhook Webhook
		var listWebhook []Webhook
		// makes variables
		// if the id is empty, put only the id in the iteration
		if webHookId != "" {
			iteration = Db.Client.Collection("webhooks").Where("ID", "==", webHookId).Documents(Db.Ctx)
		}
		// else put all 
		else { 
			iteration = Db.Client.Collection("webhooks").Documents(Db.Ctx)
		}
		
		// loops through the document
		for {
			doc, err := iteration.Next() 
			//exit once it's finished
			if err == iterator.Done {
				break 
			}
			//exit if the iteration fails
			if err != nil { 
				log.Fatalf("Failed to iterate: %v", err)
			}
			//add to the array if it got this far
			if doc != nil {
				err = doc.DataTo(&webhook)                
				listWebhook = append(listWebhook, webhook)
			}
		}

		//Encode to the json
		w.Header().Add("content-type", "application/json")
		err := json.NewEncoder(w).Encode(listWebhook)
		if err != nil {
			http.Error(w, "failed to encode webhoks", http.StatusInternalServerError)
			fmt.Println("failed to encode")
		}

		//deletes
	case http.MethodDelete:
		var deleteHook Webhook
		// finds the Id
		webhookIdentifier := r.URL.Path[23:]
		deleteHook.Id = webhookIdentifier

		//Delete it from the database
		err := Db.Delete(&deleteHook) 
		if err != nil {
			log.Fatal(err)
		}

		//POST
	case http.MethodPost:
		var postHook Webhook

		// decode from the body
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&postHook)
		if err != nil {
			fmt.Println("Failed to decode json from user in post request")
		}
		//Save it to the database
		err = Db.Save(&postHook)
		if err != nil {
			log.Fatal(err)
		}

	default:
		// default function, should never really run
		fmt.Fprintf(w, "Unknown REST method %v", r.Method)
	}
}
