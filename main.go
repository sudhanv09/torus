package main

import (
	"fmt"
	"log"
	"github.com/joho/godotenv"
	"sudhanv09/torus/metadata"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Searching for movies...")
	movies, err := metadata.SearchMovie("deadpool")
	if err != nil {
		log.Fatalf("Error searching for movies: %v", err)
	}

	fmt.Println("Movies found:", movies.TotalResults)
	for _, row := range movies.Results {
		fmt.Println(row.Title)
	}

	fmt.Println("Searching for series...")
	series, err := metadata.SearchSeries("severance")
	if err != nil {
		log.Fatalf("Error searching for series: %v", err)
	}

	fmt.Println("Series found:", series.Status, len(series.Data))
	for _, row := range series.Data {
		fmt.Println(row.Name)
	}
}
