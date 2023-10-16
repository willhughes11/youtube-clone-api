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

	// Base
	router.GET("/api/v1", getApiBaseEndpoint)

	// Videos
	router.GET("/api/v1/videos/mostPopular", getPopularYoutubeVideos)
	router.GET("/api/v1/videos/:id", getVideoById)
	router.GET("/api/v1/videos/channel/:id", getVideosByChannelId)

	// Playlists
	router.GET("/api/v1/playlists/:id", getPlaylistById)

	// Playlist Items
	router.GET("/api/v1/playlistItems/playlist/:playlist_id", getPlaylistItemsByPlaylistId)

	// Video Categories
	router.GET("/api/v1/videosCategories", getVideoCategoriesByRegionCode)

	// Channels
	router.GET("/api/v1/channels", getChannel)

	// Channel Sections
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
