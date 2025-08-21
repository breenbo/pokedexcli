package internal

import (
	"encoding/json"
	"net/http"
)

type Result struct {
	Name string
	Url  string
}

type APIResponse struct {
	Next     string
	Previous string
	Results  []Result
}

func FetchLocation() (APIResponse, error) {
	locationURL := "https://pokeapi.co/api/v2/location"
	locations := APIResponse{}

	res, err := http.Get(locationURL)
	if err != nil {
		return locations, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locations); err != nil {
		return locations, err
	}

	return locations, nil
}
