package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"time"
	"github.com/j-elliott3/pokedexcli/internal/pokecache"
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

type Client struct {
	cache *pokecache.Cache
}

func NewClient(interval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(interval),
	}
}

func (c *Client)LocationAreaGET(url string) (LocationAreasResponse, error) {
	fmt.Println("debug: LocationAreaGET called with", url)
	if data, ok := c.cache.Get(url); ok {
		fmt.Println("debug: cache hit")
		return unmarshalDataHelper[LocationAreasResponse](data)
	}
	
	fmt.Println("debug: making HTTP request")
	res, err := http.Get(url)
	fmt.Println("debug: HTTP request returned")
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	fmt.Println("debug: body read, length:", len(body))
	fmt.Println("debug: status code:", res.StatusCode)
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code %d\n", res.StatusCode)
		return LocationAreasResponse{}, fmt.Errorf("Response error")
	}
	if err != nil {
		return LocationAreasResponse{}, err
	}
	c.cache.Add(url, body)

	fmt.Println("debug: about to unmarshal")
	return unmarshalDataHelper[LocationAreasResponse](body)
}

func unmarshalDataHelper[T any](body []byte) (T, error) {
	var result T
	err := json.Unmarshal(body, &result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

type Location struct {
	Name string `json:"name"`
	PokemonEncounters []Encounter `json:"pokemon_encounters"`
}

type Encounter struct {
	Pokemon Pokemon `json:"pokemon"`
}
type Pokemon struct {
	Name string `json:"name"`
}

var locationBaseURL = "https://pokeapi.co/api/v2/location-area/"

func (c *Client)LocationGET(name string) (Location, error) {
	locationURL := locationBaseURL + name 
	if data, ok := c.cache.Get(locationURL); ok {
		return unmarshalDataHelper[Location](data)
	}
	
	res, err := http.Get(locationURL)
	if err != nil {
		return Location{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Location{}, err
	}
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code %d\n", res.StatusCode)
		return Location{}, fmt.Errorf("Response error")
	}
	c.cache.Add(locationURL, body)

	return unmarshalDataHelper[Location](body)
}