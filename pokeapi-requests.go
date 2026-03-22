package main

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)


type LocationAreasResponse struct {
    Count    int     `json:"count"`
    Next     *string `json:"next"`
    Previous *string `json:"previous"`
    Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
    URL  string `json:"url"`
}

func pokeapiLocationAreaGET(url string) (LocationAreasResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code %d\n", res.StatusCode)
		return LocationAreasResponse{}, fmt.Errorf("Response error")
	}
	if err != nil {
		return LocationAreasResponse{}, err
	}

	locationAreas := LocationAreasResponse{}
	err = json.Unmarshal(body, &locationAreas)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	return locationAreas, nil
}