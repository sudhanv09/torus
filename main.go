package main

import (
	"log"
	"sudhanv09/torus/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, continuing with environment variables")
	}

	_, err = db.InitDB("torus.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	log.Println("Database initialized successfully")

	startServer()
}

