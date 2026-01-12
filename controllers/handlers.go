package controllers

import (
	"net/http"
	"sudhanv09/torus/views/pages"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	pages.Home().Render(r.Context(), w)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Search results will appear here."))
}

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Settings and configuration."))
}