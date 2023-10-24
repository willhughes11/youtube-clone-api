package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/youtube/v3"
)

func getVideoCategoriesByRegionCode(c *gin.Context) {
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
