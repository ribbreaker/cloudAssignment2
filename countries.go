package main

import (
	"encoding/json"
	"fmt"
	"net/http"

)

//Country more like go fuck yourself
type Country struct {
	//2-letter ISO format country code
	Code string `json:"countryCode"`
	//english human-readable country name
	CountryName string `json:"countryname"`
	//Country flag
	CountryFlag string `json:"flag"`
	//Species
	Species []Species `json:"species"`
	//Species key
	SpeciesKey []string `json:"speciesKey"`
}

type CountryResponse struct{
	Collection []Country
}


//List a given number of species entries by country
func countryHandler(w http.ResponseWriter, r *http.Request) {
	country_identifier := r.URL.Path[25:]
	var limit = 20
	if r.URL.Query()["limit"] != nil {
		customLimit := r.URL.Query()["limit"][0]
		customLimitInt, err := strconv.Atoi(customLimit)
		if err == nil {
			limit = customLimitInt
		}
	}



	//restcountries
	resp, err := http.Get("https://restcountries.eu/rest/v2/alpha/" + country_identifier)
	if err != nil {
		//handle error
		fmt.Println("Error parsing request", err)
		return
	}
	defer resp.Body.Close()

	resp, err = http.Get("http://api.gbif.org/v1/occurrence/search?country=" + country_identifier + "&limit=" + strconv.Itoa(limit))
	if err != nil {
		//handle error
		fmt.Println("Error parsing request", err)
		return
	}
	defer resp.Body.Close()

	//resp.body JSON
	country := &Country{}

	fmt.Println(resp.Body, "<-- resp body")
	err = json.NewDecoder(resp.Body).Decode(country)
	if err != nil {
		//handle error
		fmt.Println("Error reading JSON data", err)
		return
	}
}
