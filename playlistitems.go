package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getPlaylistItemsByPlaylistId(c *gin.Context) {
	playlistParts := strings.Split("contentDetails,localizations,snippet,status", ",")
	playlistId := c.Param("playlist_id")

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	call := service.PlaylistItems.List(playlistParts)
	call = call.PlaylistId(playlistId)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}
