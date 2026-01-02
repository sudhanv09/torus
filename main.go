package main

import (
	"log"
	// "github.com/joho/godotenv"
	"github.com/a-h/templ"
	"net/http"
	"sudhanv09/torus/view"
)

func main() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	dashboard := view.Dashboard()
	http.Handle("/", templ.Handler(dashboard))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))

}
