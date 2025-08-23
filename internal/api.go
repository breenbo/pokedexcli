package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
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

type ExploreResponse struct {
	Pokemon_encounters []struct {
		Pokemon Result
	}
}

type PokemonRes struct {
	Name            string
	Base_experience int
	Height          int
	Weight          int
	Types           []struct {
		Type Result
	}
	Stats []struct {
		Base_stat int
		Stat      struct {
			Name string
		}
	}
}

type Pokemon struct {
	Name            string
	Base_experience int
	Height          int
	Weight          int
	Types           []string
	Stats           []map[string]int
}

func cacheOrFetch(url string) ([]byte, error) {
	duration := time.Second * 5
	cache := NewCache(duration)
	var body []byte
	// use cached value first
	if cached, ok := cache.Get(url); ok {
		body = cached
	} else {
		res, err := http.Get(url)
		if err != nil {
			return body, err
		}
		defer res.Body.Close()

		live, err := io.ReadAll(res.Body)
		if err != nil {
			return body, err
		}

		body = live
		cache.Add(url, live)

	}

	return body, nil
}

func FetchLocation(url string) (APIResponse, error) {
	locations := APIResponse{}

	body, err := cacheOrFetch(url)
	if err != nil {
		return locations, err
	}

	if err := json.Unmarshal(body, &locations); err != nil {
		return locations, err
	}

	return locations, nil
}

func FetchExplore(url string, areaName string) ([]Result, error) {
	exploreRes := ExploreResponse{}
	fullUrl := url + "/" + areaName

	body, err := cacheOrFetch(fullUrl)
	if err != nil {
		return []Result{}, err
	}

	if err := json.Unmarshal(body, &exploreRes); err != nil {
		return []Result{}, err
	}

	pokemonList := []Result{}
	for _, pokemon := range exploreRes.Pokemon_encounters {
		pokemonList = append(pokemonList, pokemon.Pokemon)
	}

	return pokemonList, nil
}

func FetchPokemon(url string, pokemonName string) (Pokemon, error) {
	fullUrl := url + "/" + pokemonName
	pokemonRes := PokemonRes{}
	pokemon := Pokemon{}

	body, err := cacheOrFetch(fullUrl)
	if err != nil {
		return pokemon, err
	}

	if err := json.Unmarshal(body, &pokemonRes); err != nil {
		return pokemon, err
	}
	pokemon.Name = pokemonRes.Name
	pokemon.Base_experience = pokemonRes.Base_experience
	pokemon.Height = pokemonRes.Height
	pokemon.Weight = pokemonRes.Weight

	pokemon.Types = make([]string, len(pokemonRes.Types))
	for i, typeInfo := range pokemonRes.Types {
		pokemon.Types[i] = typeInfo.Type.Name
	}

	pokemon.Stats = make([]map[string]int, len(pokemonRes.Stats))
	for i, statInfo := range pokemonRes.Stats {
		pokemon.Stats[i] = map[string]int{
			statInfo.Stat.Name: statInfo.Base_stat,
		}
	}

	return pokemon, nil
}
