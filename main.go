package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type apiBase struct {
	Info string `json:"info"`
}

var (
	popularVideoParts          = strings.Split("contentDetails,id,snippet,statistics,topicDetails", ",")
	channelProfilePictureParts = strings.Split("snippet", ",")
)

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Adjust this to your domain(s)
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	router.GET("/", getApiBaseEndpoint)
	router.GET("/videos/popular", getPopularYoutubeVideos)
	router.GET("/channel/:id/thumbnails", getChannelProfilePicture)

	router.Run("localhost:8080")
}

func getApiBaseEndpoint(c *gin.Context) {
	var apiBaseInfo = apiBase{
		Info: "LiveSync API",
	}

	c.IndentedJSON(http.StatusOK, apiBaseInfo)
}

func getPopularYoutubeVideos(c *gin.Context) {
	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	// Make an API request to fetch YouTube video data
	call := service.Videos.List(popularVideoParts)
	call = call.Chart("mostPopular")
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return the API response with a 200 status code
	c.IndentedJSON(http.StatusOK, response)
}

func getChannelProfilePicture(c *gin.Context) {
	channelId := c.Param("id")

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	// Make an API request to fetch YouTube video data
	call := service.Channels.List(channelProfilePictureParts)
	call = call.Id(channelId)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	channelThumbnails := map[string]interface{}{
		"etag":       response.Etag,
		"kind":       response.Kind,
		"id":         response.Items[0].Id,
		"thumbnails": response.Items[0].Snippet.Thumbnails,
	}

	// Return the API response with a 200 status code
	c.IndentedJSON(http.StatusOK, channelThumbnails)
}
