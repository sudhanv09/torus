package models

import (
	"time"
)

type TrackedItem struct {
	ID           int64     `json:"id"`
	Type         string    `json:"type"` // "movie" or "series"
	ExternalID   string    `json:"external_id"`
	Title        string    `json:"title"`
	Year         int       `json:"year"`
	PosterURL    string    `json:"poster_url"`
	BackdropPath string    `json:"backdrop_path"`
	Overview     string    `json:"overview"`
	Genres       string    `json:"genres"` // Comma-separated
	Path         string    `json:"path"`   // Disk path
	Status       string    `json:"status"` // "wanted", "monitored", "completed"
	AddedAt      time.Time `json:"added_at"`
}

type Season struct {
	ID            int64  `json:"id"`
	TrackedItemID int64  `json:"tracked_item_id"`
	SeasonNumber  int    `json:"season_number"`
	Title         string `json:"title"`
	Overview      string `json:"overview"`
	PosterPath    string `json:"poster_path"`
	Status        string `json:"status"`
}

type Episode struct {
	ID            int64     `json:"id"`
	TrackedItemID int64     `json:"tracked_item_id"`
	SeasonID      int64     `json:"season_id"`
	SeasonNumber  int       `json:"season_number"`
	EpisodeNumber int       `json:"episode_number"`
	Title         string    `json:"title"`
	AirDate       time.Time `json:"air_date"`
	Status        string    `json:"status"` // "wanted", "downloaded", "failed"
}
