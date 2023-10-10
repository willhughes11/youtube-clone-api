package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/youtube/v3"
)

func getPopularYoutubeVideos(c *gin.Context) {
	videoCategoryId := c.DefaultQuery("vcid", "0")
	regionCode := c.DefaultQuery("rc", "US")
	nextPageToken := c.DefaultQuery("npt", "")
	popularVideoParts := strings.Split("contentDetails,id,snippet,statistics,topicDetails", ",")

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	call := service.Videos.List(popularVideoParts)
	call = call.Chart("mostPopular").
		VideoCategoryId(videoCategoryId).
		RegionCode(regionCode).
		PageToken(nextPageToken)

	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Loop through the response and modify the Snippet object for each item
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Unmarshal the JSON into a Go data structure (e.g., a map)
	var data map[string]interface{}
	if err := json.Unmarshal(jsonResponse, &data); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Loop through the data and make changes as needed
	// For example, you can add a new key-value pair to each item's "snippet"
	if items, ok := data["items"].([]interface{}); ok {
		for _, item := range items {
			if itemMap, itemIsMap := item.(map[string]interface{}); itemIsMap {
				// Ensure "snippet" field exists and is a map
				snippet, snippetExists := itemMap["snippet"].(map[string]interface{})
				if !snippetExists {
					snippet = make(map[string]interface{})
					itemMap["snippet"] = snippet
				}

				channelId := snippet["channelId"]

				if channelIdStr, ok := channelId.(string); ok {
					channelThumbnails := getChannelThumbnails(service, channelIdStr)

					// Modify or add a new key-value pair to "snippet"
					snippet["channelThumbnails"] = channelThumbnails
				}

			}
		}
	}

	// Marshal the modified data back into JSON
	modifiedResponse, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling modified data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", modifiedResponse)
}

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
