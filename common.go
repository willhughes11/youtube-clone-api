package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func getGoogleApiService() (*youtube.Service, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, err
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")

	// Create a new YouTube service client
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		log.Printf("Error creating YouTube service client: %v", err)
		return nil, err
	}

	return service, nil
}
