package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	// "time"
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
	duration := time.Second * 5

	cache := NewCache(duration)
	var body []byte

	// use cached value first
	if cached, ok := cache.Get(url); ok {
		body = cached
	} else {
		res, err := http.Get(url)
		if err != nil {
			return locations, err
		}
		defer res.Body.Close()

		live, err := io.ReadAll(res.Body)
		if err != nil {
			return locations, err
		}

		body = live
		cache.Add(url, live)
	}

	if err := json.Unmarshal(body, &locations); err != nil {
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
