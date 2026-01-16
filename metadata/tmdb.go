package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sudhanv09/torus/models"
)

const baseTMDBURL = "https://api.themoviedb.org/3"

func SearchMulti(query string) (*models.TMDBResponse, error) {
	u, err := url.Parse(baseTMDBURL + "/search/multi")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("language", "en-US")
	q.Set("page", "1")
	q.Set("include_adult", "false")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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
		fmt.Println("Error multi searching:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tmdb multi search failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tmdbResponse models.TMDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		return nil, err
	}

	return &tmdbResponse, nil
}

func buildSearchURL(query string, q_type string) (string, error) {
	u, err := url.Parse(baseTMDBURL + "/search/" + q_type)
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

func SearchMovie(query string) (*models.TMDBResponse, error) {
	searchURL, err := buildSearchURL(query, "movie")
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

	var tmdbResponse models.TMDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		return nil, err
	}

	return &tmdbResponse, nil
}

func SearchTv(query string) (*models.TMDBResponse, error) {
	searchURL, err := buildSearchURL(query, "tv")
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
		fmt.Println("Error searching for tv shows:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tmdb search failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tmdbResponse models.TMDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		return nil, err
	}

	return &tmdbResponse, nil
}

// https://api.themoviedb.org/3/tv/{series_id}
func GetShowById(id int) (*models.TMDBShowDetails, error) {
	u := fmt.Sprintf("%s/tv/%d?language=en-US", baseTMDBURL, id)
	req, err := http.NewRequest("GET", u, nil)
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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tmdb get show failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var show models.TMDBShowDetails
	if err := json.NewDecoder(resp.Body).Decode(&show); err != nil {
		return nil, err
	}

	return &show, nil
}

// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}
func GetSeasonById(showId int, seasonNumber int) (*models.TMDBSeasonDetails, error) {
	u := fmt.Sprintf("%s/tv/%d/season/%d?language=en-US", baseTMDBURL, showId, seasonNumber)
	req, err := http.NewRequest("GET", u, nil)
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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tmdb get season failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var season models.TMDBSeasonDetails
	if err := json.NewDecoder(resp.Body).Decode(&season); err != nil {
		return nil, err
	}

	return &season, nil
}
