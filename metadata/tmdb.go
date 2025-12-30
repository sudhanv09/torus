package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseTMDBURL = "https://api.themoviedb.org/3"

type TMDBMovie struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	ReleaseDate  string  `json:"release_date"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	VoteAverage  float64 `json:"vote_average"`
	VoteCount    int     `json:"vote_count"`
}

type TMDBResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

func SearchMovie(query string) (*TMDBResponse, error) {
	url := fmt.Sprintf("%s/search/movie?query=%s", baseTMDBURL, query)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var tmdbResponse TMDBResponse
	json.NewDecoder(response.Body).Decode(&tmdbResponse)
	return &tmdbResponse, nil
}
