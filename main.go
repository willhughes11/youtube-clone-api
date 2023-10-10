package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	router.GET("/api/v1/videos/categories", getVideoCategoriesByRegionCode)

	router.Run()
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
