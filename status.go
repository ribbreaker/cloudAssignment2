package assignment2

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// this takes a timestamp at the start of the program
var upTime = time.Now()

func StatusHandler(w http.ResponseWriter, r *http.Request) {

	//if I get the data, it is working
	gitResp, err := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects")
	if err != nil {
		log.Fatalln(err)
	}

	//Assumes it to be ok
	dbStatus := 200
	//Doc("0") is reserved for status checks
	_, err = Db.Client.Collection("webhooks").Doc("0").Get(Db.Ctx)
	if err != nil {
		//Gives a Service Unavailable error
		dbStatus = 503
	}

	// I save them away to give to the status struct
	Statusgit := gitResp.StatusCode

	// I take a new timestamp at this point and find the difference between now and the old time
	newTime := time.Now()
	bigTime := int64(newTime.Sub(upTime) / time.Second)

	statusData := Status{Statusgit, dbStatus, "v1", bigTime}

	//Encode the json, standard stuff
	w.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(statusData)

}
