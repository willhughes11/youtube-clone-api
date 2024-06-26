package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func getPopularYoutubeVideos(c *gin.Context) {
	videoCategoryId := c.DefaultQuery("vcid", "0")
	regionCode := c.DefaultQuery("rc", "US")
	getChannelProfilePicture := c.DefaultQuery("gcpp", "true")
	nextPageToken := c.Query("npt")
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

	if getChannelProfilePicture == "true" {
		// Loop through the response and modify the Snippet object for each item
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		data := processItemsConcurrently(jsonResponse, service, false, "")

		// Marshal the modified data back into JSON
		modifiedResponse, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshaling modified data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.Data(http.StatusOK, "application/json; charset=utf-8", modifiedResponse)
	} else {
		c.IndentedJSON(http.StatusOK, response)
	}
}

func getVideoById(c *gin.Context) {
	videoParts := strings.Split("contentDetails,id,snippet,statistics,topicDetails", ",")
	videoId := c.Param("id")

	service, err := getGoogleApiService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	call := service.Videos.List(videoParts)
	call = call.Id(videoId)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making API call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func getVideosByChannelId(c *gin.Context) {
	videoSearchParts := strings.Split("snippet", ",")
	channelId := c.Param("id")
	order := c.DefaultQuery("order", "date")
	getVideoItemObject := c.DefaultQuery("gvio", "true")
	getFakeData := c.DefaultQuery("gfd", "false")
	nextPageToken := c.DefaultQuery("npt", "")

	if getFakeData == "false" {
		service, err := getGoogleApiService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}

		call := service.Search.List(videoSearchParts)
		call = call.ChannelId(channelId).Type("video").Order(order).VideoLicense("youtube")
		response, err := call.Do()
		if err != nil {
			log.Printf("Error making API call: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		if getVideoItemObject == "true" {
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			data := processItemsConcurrently(jsonResponse, service, true, nextPageToken)

			// Marshal the modified data back into JSON
			modifiedResponse, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshaling modified data: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			c.Data(http.StatusOK, "application/json; charset=utf-8", modifiedResponse)

		} else {
			c.IndentedJSON(http.StatusOK, response)
		}
	} else {
		data, err := os.ReadFile("json/videos-by-channel-id.json")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to read data from the JSON file"})
			return
		}

		var items map[string]interface{}
		if err := json.Unmarshal(data, &items); err != nil {
			c.JSON(500, gin.H{"error": "Failed to unmarshal JSON"})
			return
		}

		time.Sleep(1 * time.Second)

		c.JSON(200, items)
	}
}
