package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
