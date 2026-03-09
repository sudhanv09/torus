package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"sudhanv09/torus/metadata"
	"sudhanv09/torus/models"
	"sudhanv09/torus/views/components"
)

func MediaDetailsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	mediaType := r.URL.Query().Get("type")

	id, err := strconv.Atoi(idStr)
	if err != nil || (mediaType != "movie" && mediaType != "tv") {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}

	var details *models.MediaDetails
	switch mediaType {
	case "movie":
		movie, err := metadata.GetMovieById(id)
		if err != nil {
			log.Printf("error fetching movie %d: %v", id, err)
			http.Error(w, "failed to fetch movie details", http.StatusInternalServerError)
			return
		}
		details = movieToDetails(movie)

	case "tv":
		show, err := metadata.GetShowById(id)
		if err != nil {
			log.Printf("error fetching show %d: %v", id, err)
			http.Error(w, "failed to fetch show details", http.StatusInternalServerError)
			return
		}
		details = showToDetails(show)
	}

	alreadyTracked, _ := models.IsTracked(mediaType, strconv.Itoa(id))
	components.MediaDialogContent(details, alreadyTracked).Render(r.Context(), w)
}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	var signals struct {
		TmdbId    float64 `json:"tmdbId"`
		MediaType string  `json:"mediaType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&signals); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id := int(signals.TmdbId)
	mediaType := signals.MediaType

	if id == 0 || (mediaType != "movie" && mediaType != "tv") {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}

	externalID := strconv.Itoa(id)
	already, _ := models.IsTracked(mediaType, externalID)
	if already {
		w.Write([]byte(`<div id="dialog-content"></div>`))
		return
	}

	switch mediaType {
	case "movie":
		if err := trackMovie(id, externalID); err != nil {
			log.Printf("error tracking movie %d: %v", id, err)
			http.Error(w, "failed to track", http.StatusInternalServerError)
			return
		}

	case "tv":
		if err := trackShow(id, externalID); err != nil {
			log.Printf("error tracking show %d: %v", id, err)
			http.Error(w, "failed to track", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte(`<div id="dialog-content"></div>`))
}

// --- helpers ---

func trackMovie(id int, externalID string) error {
	movie, err := metadata.GetMovieById(id)
	if err != nil {
		return err
	}

	genres := genreNames(movie.Genres)
	year := 0
	if len(movie.ReleaseDate) >= 4 {
		year, _ = strconv.Atoi(movie.ReleaseDate[:4])
	}

	item := &models.TrackedItem{
		Type:         "movie",
		ExternalID:   externalID,
		Title:        movie.Title,
		Year:         year,
		PosterURL:    downloadPoster(movie.PosterPath, id),
		BackdropPath: movie.BackdropPath,
		Overview:     movie.Overview,
		Genres:       genres,
		Status:       "wanted",
	}
	return models.AddTrackedItem(item)
}

func trackShow(id int, externalID string) error {
	show, err := metadata.GetShowById(id)
	if err != nil {
		return err
	}

	genres := genreNames(show.Genres)
	year := 0
	if len(show.FirstAirDate) >= 4 {
		year, _ = strconv.Atoi(show.FirstAirDate[:4])
	}

	item := &models.TrackedItem{
		Type:         "tv",
		ExternalID:   externalID,
		Title:        show.Name,
		Year:         year,
		PosterURL:    downloadPoster(show.PosterPath, id),
		BackdropPath: show.BackdropPath,
		Overview:     show.Overview,
		Genres:       genres,
		Status:       "wanted",
	}
	if err := models.AddTrackedItem(item); err != nil {
		return err
	}

	for _, s := range show.Seasons {
		if s.SeasonNumber == 0 {
			continue // skip specials
		}
		season := &models.Season{
			TrackedItemID: item.ID,
			SeasonNumber:  s.SeasonNumber,
			Title:         s.Name,
			Overview:      s.Overview,
			PosterPath:    s.PosterPath,
			Status:        "wanted",
		}
		if err := models.AddSeason(season); err != nil {
			log.Printf("warn: failed to save season %d for show %d: %v", s.SeasonNumber, id, err)
		}
	}

	return nil
}

func movieToDetails(m *models.TMDBMovieDetails) *models.MediaDetails {
	year := ""
	if len(m.ReleaseDate) >= 4 {
		year = m.ReleaseDate[:4]
	}
	return &models.MediaDetails{
		ID:           m.ID,
		Type:         "movie",
		Title:        m.Title,
		Year:         year,
		Overview:     m.Overview,
		PosterPath:   m.PosterPath,
		BackdropPath: m.BackdropPath,
		VoteAverage:  m.VoteAverage,
		Genres:       genreNames(m.Genres),
		Runtime:      m.Runtime,
	}
}

func showToDetails(s *models.TMDBShowDetails) *models.MediaDetails {
	year := ""
	if len(s.FirstAirDate) >= 4 {
		year = s.FirstAirDate[:4]
	}
	return &models.MediaDetails{
		ID:           s.ID,
		Type:         "tv",
		Title:        s.Name,
		Year:         year,
		Overview:     s.Overview,
		PosterPath:   s.PosterPath,
		BackdropPath: s.BackdropPath,
		VoteAverage:  s.VoteAverage,
		Genres:       genreNames(s.Genres),
		NumSeasons:   s.NumberOfSeasons,
	}
}

func downloadPoster(tmdbPath string, id int) string {
	if tmdbPath == "" {
		return ""
	}
	url := fmt.Sprintf("https://image.tmdb.org/t/p/w342%s", tmdbPath)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("warn: failed to download poster for %d: %v", id, err)
		return ""
	}
	defer resp.Body.Close()

	localPath := fmt.Sprintf("static/posters/%d.jpg", id)
	f, err := os.Create(localPath)
	if err != nil {
		log.Printf("warn: failed to create poster file for %d: %v", id, err)
		return ""
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		log.Printf("warn: failed to write poster for %d: %v", id, err)
		return ""
	}
	return "/static/posters/" + fmt.Sprintf("%d.jpg", id)
}

func genreNames(genres []models.TMDBGenre) string {
	names := make([]string, 0, len(genres))
	for _, g := range genres {
		names = append(names, g.Name)
	}
	return strings.Join(names, ", ")
}
