package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @BasePath /channels

// PingExample godoc
// @Summary Get Channels
// @Schemes
// @Description Base Endpoint
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} YouTube Clone API
// @Router /channels [get]
func getChannel(c *gin.Context) {
	channelParts := strings.Split("id,contentDetails,id,snippet,statistics,topicDetails,status,brandingSettings,localizations", ",")
	channelId := c.Query("id")
	username := c.Query("uname")

	if len(channelId) > 0 && len(username) > 0 || len(channelId) == 0 && len(username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	call := service.Channels.List(channelParts)
	call = call.Id(channelId)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}
