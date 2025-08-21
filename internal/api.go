package internal

import (
	"encoding/json"
	"fmt"
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

func FetchLocation(url string) (APIResponse, error) {
	locations := APIResponse{}

	res, err := http.Get(url)
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

func PrintResults(results []Result) {
	// print the locations from slice
	locationMsg := "\n"
	for _, loc := range results {
		locationMsg += fmt.Sprintf("%s\n", loc.Name)
	}
	locationMsg += "\n"

	fmt.Print(locationMsg)
}
