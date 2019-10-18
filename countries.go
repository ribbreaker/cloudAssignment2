package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

//Country golint wont leave me alone
type Country struct {
	//2-letter ISO format country code
	Code string `json:"alpha2Code"`
	//english human-readable country name
	CountryName string `json:"name"`
	//Country flag
	CountryFlag string `json:"flag"`
	//Species
	Species []Species `json:"species"`
	//Species key
	SpeciesKey []int `json:"speciesKey"`
}

//CountryResponse golint wont stop torturing me
type CountryResponse struct {
	Species    Species
	SpeciesKey int
}

//List a given number of species entries by country
func countryHandler(w http.ResponseWriter, r *http.Request) {
	countryIdentifier := r.URL.Path[25:]

	// Country codes are two characters long, if not Error
	/*if len(countryCode) != 2 {
		http.Error(w, "Wrong country code used", http.StatusBadRequest)
	}*/

	
	var limit = 20
	if r.URL.Query()["limit"] != nil {
		customLimit := r.URL.Query()["limit"][0]
		customLimitInt, err := strconv.Atoi(customLimit)
		if err == nil {
			limit = customLimitInt
		}
	}

	//restcountries
	resp, err := http.Get("https://restcountries.eu/rest/v2/alpha/" + countryIdentifier)
	if err != nil {
		//handle error
		fmt.Println("Error parsing request", err)
		return
	}
	defer resp.Body.Close()

	//resp.body JSON
	country := &Country{}
	err = json.NewDecoder(resp.Body).Decode(country)

	resp, err = http.Get("http://api.gbif.org/v1/occurrence/search?country=" + countryIdentifier + "&limit=" + strconv.Itoa(limit))
	if err != nil {
		//handle error
		fmt.Println("Error parsing request", err)
		return
	}
	defer resp.Body.Close()

	//Secondary struct
	countryResponse := &CountryResponse{}
	err = json.NewDecoder(resp.Body).Decode(countryResponse)

	//Loops through all of them and assigns species/specieskey

	//(use map for not getting duplicates)
	/*	for x := 0; x <= limit; x++ {
		//Gets the species key
		if countryResponse.SpeciesKey[x] != "" {
			country.SpeciesKey[x] = countryResponse.SpeciesKey[x]

			//if there's a specieskey, got get that species
			resp, err := http.Get("http://api.gbif.org/v1/species/" + countryResponse.SpeciesKey[x])
			if err != nil {
				log.Fatalln(err)
			}
			//set the Country's species to be the species just pulled out
			var mySpecies Species
			err = json.NewDecoder(resp.Body).Decode(&mySpecies)
			country.Species[x] = mySpecies
		}

	}*/
//	Var m = make(map[Species]int)

	err = json.NewDecoder(resp.Body).Decode(country)
	if err != nil {
		//handle error
		fmt.Println("Error reading JSON data", err)
		return
	}

	w.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(country)
}
/*

func countryHandler (w http.ResponseWriter, r *http.Request) {

	// The URL is 25 + key long, so I skip the 25 first keys
	countryCode := r.URL.Path[25:]

	// Pull the country information from the website
	respCountry, err := http.Get(fmt.Sprintf("https://restcountries.eu/rest/v2/alpha/%s", countryCode))
	if err != nil {
		log.Fatalln(err)
	}

	var countryData Country
	var resultData Result

	//  Since I don't actually care about the error code, I use _
	_ = json.NewDecoder(respCountry.Body).Decode(&countryData)

	fmt.Printf("decode\n")

					// Since the default limit for lookups is 20, I figured the best way to include
	limit := 20		// both a limit and not, is to set the limit to 20. That way I can do both limit-checks and
	offset := 0		// no-limit-checks.

	fmt.Printf("for limit nr -1: %d\n", limit)
	if r.URL.Query()["limit"] != nil {		// see if limit is in the URL

		fmt.Printf("if limit\n")

		customLimit := r.URL.Query()["limit"][0]			// finds the string of the "number"
		customLimitInt, err := strconv.Atoi(customLimit)	// converts from string to int

		if err == nil {
			limit = customLimitInt							// if all goes well, we now have a good limit
		} else {
			log.Fatalln(err)
		}
	}	//and if the user doesn't write a limit, we use limit = 20, with the user noticing no difference
	fmt.Printf("for limit nr 0: %d\n", limit)

	for ; limit > 0; {
		fmt.Printf("for limit nr 1: %d\n", limit)
		respSpecies, err := http.Get(fmt.Sprintf("http://api.gbif.org/v1/occurrence/search?country=%s&limit=%d&offset=%d", countryCode, limit, offset))
		if err != nil {
			log.Fatalln(err)
		}

		_ = json.NewDecoder(respSpecies.Body).Decode(&resultData)


		limit -= 300
		offset += 300
		exist := false
		fmt.Printf("for limit nr 2: %d\n", limit)
		fmt.Printf("for results: %d\n", len(resultData.Results))
		for i := 0 ; i < len(resultData.Results) ; i++ {

			for j := 0; j < len(countryData.SpeciesKey); j++ {
				if resultData.Results[i].Key == countryData.SpeciesKey[j]{
					fmt.Printf("nr %d, %d == %d\n", j, resultData.Results[j].Key, countryData.SpeciesKey[j])
					exist = true
				}
			}


		if exist == false {
			countryData.SpeciesKey = append(countryData.SpeciesKey, resultData.Results[i].Key)
			countryData.Species = append(countryData.Species, resultData.Results[i].Name)
			fmt.Printf("for append %d\n", i)
		}
		exist = false


		}

	}
	fmt.Printf("for append -1: %d\n", len(resultData.Results))


	for _, i := range countryData.SpeciesKey {
		if resultData.Results[i].Key == countryData.SpeciesKey[i]{
		}
	}


	// This prints the json struct to the website
	fmt.Printf("slutt encode\n")
	w.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(countryData)
	fmt.Printf("%s\n", countryData.Code)
}
*/