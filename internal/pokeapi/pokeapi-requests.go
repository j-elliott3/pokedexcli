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
		cache: pokecache.NewCache(interval)
	}
}

func (c *Client)LocationAreaGET(url string) (LocationAreasResponse, error) {
	if data, ok := c.cache.Get(url); ok {
		return unmarshalDataHelper(data)
	}
	
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
	c.cache.Add(url, body)

	return unmarshalDataHelper(body)
}

func unmarshalDataHelper(body []byte) (LocationAreasResponse, error) {
	locationAreas := LocationAreasResponse{}
	err := json.Unmarshal(body, &locationAreas)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	return locationAreas, nil
}