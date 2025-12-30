package main

import (
	"fmt"
	"log"
	"sudhanv09/torus/scrapers"
)

func main() {
	rows, err := scrapers.Search("deadpool")
	if err != nil {
		log.Fatalf("Error searching for torrents: %v", err)
	}

	for _, row := range rows {
		fmt.Println(row.Title)
	}
}
