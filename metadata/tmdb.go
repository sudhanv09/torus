package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const baseTMDBURL = "https://api.themoviedb.org/3"

type TMDBMovie struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	ID               int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	Title            string  `json:"title"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type TMDBResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

func buildSearchURL(query string) (string, error) {
	u, err := url.Parse(baseTMDBURL + "/search/movie")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("language", "en-US")
	q.Set("page", "1")
	q.Set("include_adult", "false")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func SearchMovie(query string) (*TMDBResponse, error) {
	searchURL, err := buildSearchURL(query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	key := os.Getenv("TMDB_API")
	if key == "" {
		return nil, errors.New("TMDB_API is not set")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error searching for movies:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tmdb search failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tmdbResponse TMDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		return nil, err
	}

	return &tmdbResponse, nil
}
