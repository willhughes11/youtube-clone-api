package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/youtube-clone-api/docs"
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
		eg := v1.Group("/")
		{
			eg.GET("/", getApiBaseEndpoint)
		}
	}

	// // Base
	// router.GET("/api/v1", getApiBaseEndpoint)

	// // Videos
	// router.GET("/api/v1/videos/mostPopular", getPopularYoutubeVideos)
	// router.GET("/api/v1/videos/:id", getVideoById)
	// router.GET("/api/v1/videos/channel/:id", getVideosByChannelId)

	// // Playlists
	// router.GET("/api/v1/playlists/:id", getPlaylistById)

	// // Playlist Items
	// router.GET("/api/v1/playlistItems/playlist/:playlist_id", getPlaylistItemsByPlaylistId)

	// // Video Categories
	// router.GET("/api/v1/videosCategories", getVideoCategoriesByRegionCode)

	// // Channels
	// router.GET("/api/v1/channels", getChannel)

	// // Channel Sections
	// router.GET("/api/v1/channelSections/channel/:id", getChannelSectionsByChannelId)
	// router.GET("/api/v1/channelSections/:id", getChannelSectionsById)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run()
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func getApiBaseEndpoint(c *gin.Context) {
	type apiBase struct {
		Info string `json:"info"`
	}

	var apiBaseInfo = apiBase{
		Info: "YouTube Clone API",
	}

	c.IndentedJSON(http.StatusOK, apiBaseInfo)
}
