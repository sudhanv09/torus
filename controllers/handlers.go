package controllers

import (
	"net/http"
	"sudhanv09/torus/views/pages"
	"sudhanv09/torus/metadata"
	"sudhanv09/torus/models"
	"sudhanv09/torus/views/components"
	"encoding/json"
	"fmt"
	"log"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	pages.Home().Render(r.Context(), w)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	pages.Search().Render(r.Context(), w)
}

func SearchResultsHandler(w http.ResponseWriter, r *http.Request) {
	datastarparam := r.URL.Query().Get("datastar")
	var signals map[string]interface{}
	if err := json.Unmarshal([]byte(datastarparam), &signals); err != nil {
		log.Printf("Error unmarshalling signals: %v", err)
		return
	}

	query, _ := signals["query"].(string)

	if query == "" {
		w.Write([]byte(`<div id="search-results" class="empty-state"><p>Start typing to search for movies or tv shows</p></div>`))
		return
	}

	var results *models.TMDBResponse
	var err error

	results, err = metadata.SearchMulti(query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`<div id="search-results" class="empty-state"><p>Error searching: %v</p></div>`, err)))
		return
	}

	if results == nil || len(results.Results) == 0 {
		w.Write([]byte(`<div id="search-results" class="empty-state"><p>No results found</p></div>`))
		return
	}

	components.SearchResults(results.Results).Render(r.Context(), w)
}

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Settings and configuration."))
}