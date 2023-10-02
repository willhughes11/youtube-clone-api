package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type apiBase struct {
	Info string `json:"info"`
}

var (
	videoParts = strings.Split("contentDetails,id,snippet,statistics,topicDetails", ",")
)

func main() {
	router := gin.Default()
	router.GET("/", getApiBaseEndpoint)
	router.GET("/videos/popular", getPopularYoutubeVideos)

	router.Run("localhost:8080")
}

func getApiBaseEndpoint(c *gin.Context) {
	var apiBaseInfo = apiBase{
		Info: "LiveSync API",
	}

	c.IndentedJSON(http.StatusOK, apiBaseInfo)
}

func getPopularYoutubeVideos(c *gin.Context) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")

	// Create a new YouTube service client
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating YouTube service client: %v", err)
	}

	// Make an API request to fetch YouTube video data
	call := service.Videos.List(videoParts)
	call = call.Chart("mostPopular")
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call: %v", err)
	}

	// Print the API response
	fmt.Printf("%+v\n", response)
	c.IndentedJSON(http.StatusOK, response)
}
