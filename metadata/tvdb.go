package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseTVDBURL = "https://api.thetvdb.com"

type TVDBSeries struct {
	Aliases         []string `json:"aliases"`
	Companies       []string `json:"companies"`
	CompanyType     string   `json:"companyType"`
	Country         string   `json:"country"`
	Director        string   `json:"director"`
	FirstAirTime    string   `json:"first_air_time"`
	Genres          []string `json:"genres"`
	ID              string   `json:"id"`
	ImageUrl        string   `json:"imageUrl"`
	Name            string   `json:"name"`
	IsOfficial      bool     `json:"isOfficial"`
	Network         string   `json:"network"`
	ObjectID        string   `json:"objectID"`
	Overview        string   `json:"overview"`
	Poster          string   `json:"poster"`
	Posters         []string `json:"posters"`
	PrimaryLanguage string   `json:"primary_language"`
	Status          string   `json:"status"`
	Slug            string   `json:"slug"`
	Studios         []string `json:"studios"`
	Title           string   `json:"title"`
	Thumbnail       string   `json:"thumbnail"`
	TVDBID          string   `json:"tvdb_id"`
	Type            string   `json:"type"`
	Year            string   `json:"year"`
}

type TVDBResponse struct {
	Data   []TVDBSeries `json:"data"`
	Status string       `json:"status"`
	Links  struct {
		Self       string `json:"self"`
		Next       int    `json:"next"`
		Previous   int    `json:"prev"`
		TotalItems int    `json:"total_items"`
		PageSize   int    `json:"page_size"`
	}
}

type TVDBLoginResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Status string `json:"status"`
}

func login() (*TVDBLoginResponse, error) {
	url := fmt.Sprintf("%s/login", baseTVDBURL)
	reqBody := map[string]string{
		"apikey": "",
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var tvdbLoginResponse TVDBLoginResponse
	json.NewDecoder(response.Body).Decode(&tvdbLoginResponse)
	return &tvdbLoginResponse, nil
}

func SearchSeries(query string) (*TVDBResponse, error) {
	url := fmt.Sprintf("%s/search/series?query=%s", baseTVDBURL, query)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var tvdbResponse TVDBResponse
	json.NewDecoder(response.Body).Decode(&tvdbResponse)
	return &tvdbResponse, nil
}
