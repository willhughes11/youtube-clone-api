package main

import (
	"net/http"

	docs "youtubeclone/ginapi/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// "github.com/swaggo/swag/example/basic/docs"
)

func main() {
	gin.ForceConsoleColor()
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	r.Use(cors.New(config))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		// Base
		base := v1.Group("/")
		{
			base.GET("/", getApiBaseEndpoint)
		}

		// Videos
		videos := v1.Group("videos")
		{
			videos.GET("mostPopular", getPopularYoutubeVideos)
			videos.GET(":id", getVideoById)
			videos.GET("channel/:id", getVideosByChannelId)
		}

		// Playlists
		playlists := v1.Group("playlists")
		{
			playlists.GET(":id", getPlaylistById)
		}

		// Playlist Items
		playlistsItems := v1.Group("playlistItems")
		{
			playlistsItems.GET("playlist/:playlist_id", getPlaylistItemsByPlaylistId)
		}

		// Video Categories
		videoCategories := v1.Group("videosCategories")
		{
			videoCategories.GET("/", getVideoCategoriesByRegionCode)
		}

		// Channels
		channels := v1.Group("channels")
		{
			channels.GET("/", getChannel)
		}

		// Channel Sections
		channelSections := v1.Group("channelSections")
		{
			channelSections.GET("channel/:id", getChannelSectionsByChannelId)
			channelSections.GET(":id", getChannelSectionsById)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run()
}

// @BasePath /

// PingExample godoc
// @Summary Base Endpoint
// @Schemes
// @Description Base Endpoint
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} YouTube Clone API
// @Router / [get]
func getApiBaseEndpoint(c *gin.Context) {
	type apiBase struct {
		Info string `json:"info"`
	}

	var apiBaseInfo = apiBase{
		Info: "YouTube Clone API",
	}

	c.IndentedJSON(http.StatusOK, apiBaseInfo)
}
