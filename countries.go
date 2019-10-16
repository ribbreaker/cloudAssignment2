package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

//CountriesStorage something something link wont stop bothering me
type CountriesStorage interface {
}

//Country more like go fuck yourself
type Country struct {
	Code        string    `json:"countryCode"` //2-letter ISO format country code
	CountryName string    `json:"country"`     //english human-readable country name
	CountryFlag string    `json:"countryflag"`
	Species     []Species `json:"species"`
	SpeciesKey  []string  `json:"specieskey"`
}

type allCountries []Country

//KeysResponse should have a comment or something
type KeysResponse struct {
	Collection []Country
}

//List a given number of species entries by country
func countryHandler(w http.ResponseWriter, r *http.Request) {
	//country/identifier/limit,

	url := r.URL.Path

	//identifier
	country_identifier := url[25:]
	if len(country_identifier) == 0 {
		http.Error(w, "No country identifier found", http.StatusBadRequest)
	}

	//limit

	limit := 20
	if r.URL.Query()["limit"] != nil {
		customLimit := r.URL.Query()["limit"][0]
		customLimitInt, err := strconv.Atoi(customLimit)
		if err == nil {
			limit = customLimitInt
		}
	}

	resp, err := http.Get("htpp://api.gbif.org/v1/search/" + country_identifier + strconv.Itoa(limit))

	resp, err = http.Get(url)
	if err != nil {
		//handle error
		fmt.Println("Error parsing request", err)
		return
	}
	defer resp.Body.Close()

	//resp.body JSON
	country := &Country{}
	err = json.NewDecoder(resp.Body).Decode(country)
	if err != nil {
		//handle error
		fmt.Println("Error reading JSON data", err)
		return
	}
}
