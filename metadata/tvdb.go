package metadata

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const baseTVDBURL = "https://api4.thetvdb.com/v4"

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
		Previous   string    `json:"prev"`
		Self       string `json:"self"`
		Next       string    `json:"next"`
		TotalItems int    `json:"total_items"`
		PageSize   int    `json:"page_size"`
	} `json:"links"`
}

type TVDBLoginResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Status string `json:"status"`
}

func login() (*TVDBLoginResponse, error) {
	url := fmt.Sprintf("%s/login", baseTVDBURL)
	key := os.Getenv("TVDB_API")
	if key == "" {
		return nil, errors.New("TVDB API key is not set")
	}
	reqBody := map[string]string{
		"apikey": key,
		"pin":    "",
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

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("tvdb login failed: status %d, body: %s", response.StatusCode, string(body))
	}

	var tvdbLoginResponse TVDBLoginResponse
	if err := json.NewDecoder(response.Body).Decode(&tvdbLoginResponse); err != nil {
		return nil, err
	}

	if tvdbLoginResponse.Data.Token == "" {
		return nil, errors.New("tvdb login failed: missing token in response")
	}

	return &tvdbLoginResponse, nil
}

func buildTVDBURL(query string) (string, error) {
	u, err := url.Parse(baseTVDBURL + "/search")
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("query", query)
	q.Set("type", "series")
	q.Set("language", "eng")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func SearchSeries(query string) (*TVDBResponse, error) {
	searchURL, err := buildTVDBURL(query)
	if err != nil {
		return nil, err
	}
	loginResp, err := login()
	if err != nil {
		fmt.Println("Error logging in to TVDB:", err)
		return nil, err
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginResp.Data.Token))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("tvdb search failed: status %d, body: %s", response.StatusCode, string(body))
	}

	var tvdbResponse TVDBResponse
	if err := json.NewDecoder(response.Body).Decode(&tvdbResponse); err != nil {
		return nil, err
	}
	return &tvdbResponse, nil
}
