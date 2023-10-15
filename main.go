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

	router.GET("/api/v1/videos/mostPopular", getPopularYoutubeVideos)
	router.GET("/api/v1/videos/channel/:id", getVideosByChannelId)

	router.GET("/api/v1/videosCategories", getVideoCategoriesByRegionCode)

	router.GET("/api/v1/channels", getChannel)

	router.GET("/api/v1/channelSections/channel/:id", getChannelSectionsByChannelId)
	router.GET("/api/v1/channelSections/:id", getChannelSectionsById)

	router.Run()
}

func getApiBaseEndpoint(c *gin.Context) {
	type apiBase struct {
		Info string `json:"info"`
	}

	var apiBaseInfo = apiBase{
		Info: "YouTube Clone API",
	}

	c.IndentedJSON(http.StatusOK, apiBaseInfo)
}
