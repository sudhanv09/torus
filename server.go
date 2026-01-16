package main

import (
	"log"
	"net/http"
	"sudhanv09/torus/controllers"
)

func startServer() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", controllers.HomeHandler)
	mux.HandleFunc("/search", controllers.SearchHandler)
	mux.HandleFunc("/search/results", controllers.SearchResultsHandler)
	mux.HandleFunc("/settings", controllers.SettingsHandler)

	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
