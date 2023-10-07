package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/youtube/v3"
)

func main() {
	gin.ForceConsoleColor()
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Adjust this to your domain(s)
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	router.GET("/api/v1", getApiBaseEndpoint)
	router.GET("/api/v1/videos/most-popular", getPopularYoutubeVideos)
	router.GET("/api/v1/channel/:id/thumbnails", getChannelProfileThumbnails)
	router.GET("/api/v1/videos/categories", getPopularVideoCategoriesByRegionCode)

	router.Run("localhost:8080")
}

func getApiBaseEndpoint(c *gin.Context) {
	type apiBase struct {
		Info string `json:"info"`
	}

	var apiBaseInfo = apiBase{
		Info: "LiveSync API",
	}

	c.IndentedJSON(http.StatusOK, apiBaseInfo)
}

func getPopularYoutubeVideos(c *gin.Context) {
	videoCategoryId := c.DefaultQuery("vcid", "0")
	popularVideoParts := strings.Split("contentDetails,id,snippet,statistics,topicDetails", ",")
	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	call := service.Videos.List(popularVideoParts)
	call = call.Chart("mostPopular").
		VideoCategoryId(videoCategoryId)

	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func getChannelProfileThumbnails(c *gin.Context) {
	channelProfilePictureParts := strings.Split("snippet", ",")
	channelId := c.Param("id")

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

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

	c.IndentedJSON(http.StatusOK, channelThumbnails)
}

func getPopularVideoCategoriesByRegionCode(c *gin.Context) {
	videoCategoriesByRegionParts := strings.Split("snippet", ",")
	regionCode := c.DefaultQuery("rc", "US")

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	call := service.VideoCategories.List(videoCategoriesByRegionParts)
	call = call.RegionCode(regionCode)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	defaultVideoCategoryItem := &youtube.VideoCategory{
		Etag: "default",
		Kind: response.Items[0].Kind,
		Id:   "0",
		Snippet: &youtube.VideoCategorySnippet{
			Assignable: true,
			ChannelId:  response.Items[0].Snippet.ChannelId,
			Title:      "All",
		},
	}

	response.Items = append([]*youtube.VideoCategory{defaultVideoCategoryItem}, response.Items...)

	c.IndentedJSON(http.StatusOK, response)
}
