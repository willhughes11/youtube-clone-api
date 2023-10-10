package main

import (
	"log"
	"strings"

	"google.golang.org/api/youtube/v3"
)

func getChannelThumbnails(service *youtube.Service, channelId string) *youtube.ThumbnailDetails {
	channelProfilePictureParts := strings.Split("snippet", ",")

	call := service.Channels.List(channelProfilePictureParts)
	call = call.Id(channelId)
	channelResponse, err := call.Do()

	if err != nil {
		log.Printf("Error making API call: %v", err)
		return nil // Handle error appropriately in your code
	}

	return channelResponse.Items[0].Snippet.Thumbnails
}
