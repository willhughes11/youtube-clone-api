package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getChannelSectionsByChannelId(c *gin.Context) {
	channelSectionsPart := strings.Split("id,contentDetails,snippet", ",")
	channelId := c.Param("id")

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	call := service.ChannelSections.List(channelSectionsPart)
	call = call.ChannelId(channelId)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func getChannelSectionsContentByChannelId(c *gin.Context) {
	channelSectionsPart := strings.Split("id,contentDetails,snippet", ",")
	channelId := c.Param("id")

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	call := service.ChannelSections.List(channelSectionsPart)
	call = call.ChannelId(channelId)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func getChannelSectionsById(c *gin.Context) {
	channelSectionsPart := strings.Split("id,contentDetails,snippet", ",")
	id := c.Param("id")

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	call := service.ChannelSections.List(channelSectionsPart)
	call = call.Id(id)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}
