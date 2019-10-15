package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//CountriesStorage something something link wont stop bothering me
type CountriesStorage interface {
}

//Country more like go fuck yourself
type Country struct {
	Code        string    `json:"code"`        //2-letter ISO format country code
	CountryName string    `json:"countryname"` //english human-readable country name
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

	//limit

	/*limit := 20
	if r.URL.Query()["limit"] != nil {
		customLimit := r.URL.Query()["limit"][0]
		customLimitInt, err := strconv.Atoi(customLimit)
		if err == nil {
			limit = customLimitInt
		}
	}
	*/
	url := r.RequestURI

	resp, err := http.Get(url)
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
