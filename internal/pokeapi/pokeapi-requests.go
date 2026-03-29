package pokeapi

import (
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
	return fetchAndCache[LocationAreasResponse](c, url)
}

type Location struct {
	Name string `json:"name"`
	PokemonEncounters []Encounter `json:"pokemon_encounters"`
}

type Encounter struct {
	Pokemon Pokemon `json:"pokemon"`
}
type Pokemon struct {
	Name 		string `json:"name"`
	BaseExp 	int `json:"base_experience"`
	Height 		int	`json:"height"`
	Weight 		int `json:"weight"`
	Stats		[]Stat `json:"stats"`
	Types		[]PokemonType `json:"types"`
}
type StatDetail struct {
	Name string `json:"name"`
}

type Stat struct {
	BaseStat 	int `json:"base_stat"`
	Stat		StatDetail `json:"stat"`
	
}

type TypeDetail struct {
	Name string `json:"name"`
}

type PokemonType struct {
	Type 	TypeDetail `json:"type"`
}

var locationBaseURL = "https://pokeapi.co/api/v2/location-area/"

func (c *Client)LocationGET(name string) (Location, error) {
	locationURL := locationBaseURL + name 
	return fetchAndCache[Location](c, locationURL)
}

var pokemonURL = "https://pokeapi.co/api/v2/pokemon/"

func (c *Client)pokemonGET(name string) (Pokemon, error) {
	url := pokemonURL + name + "/"
	return fetchAndCache[Pokemon](c, url)
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

func fetchAndCache[T any](c *Client, url string) (T, error) {
    if data, ok := c.cache.Get(url); ok {
        return unmarshalDataHelper[T](data)
    }

    res, err := http.Get(url)
    if err != nil {
        var zero T
        return zero, err
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        var zero T
        return zero, err
    }

    c.cache.Add(url, body)
    return unmarshalDataHelper[T](body)
}