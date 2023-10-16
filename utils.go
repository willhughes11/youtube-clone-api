package main

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	"google.golang.org/api/youtube/v3"
)

func getChannelThumbnails(service *youtube.Service, channelId string) *youtube.ThumbnailDetails {
	channelProfilePictureParts := strings.Split("snippet", ",")

	call := service.Channels.List(channelProfilePictureParts)
	call = call.Id(channelId)
	channelResponse, err := call.Do()

	if err != nil {
		log.Printf("Error making API call: %v", err)
		return nil
	}

	return channelResponse.Items[0].Snippet.Thumbnails
}

func getVideo(service *youtube.Service, videoId string) *youtube.VideoListResponse {
	videoParts := strings.Split("contentDetails,id,snippet,statistics,topicDetails", ",")

	call := service.Videos.List(videoParts)
	call = call.Id(videoId)
	response, err := call.Do()

	if err != nil {
		log.Printf("Error making API call: %v", err)
		return nil
	}

	return response
}

func processItemsConcurrently(jsonResponse []byte, service *youtube.Service, replaceItemObj bool) map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonResponse, &data); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return nil
	}

	resultChannel := make(chan map[string]interface{})

	processedItems := []map[string]interface{}{}

	if items, ok := data["items"].([]interface{}); ok {
		var wg sync.WaitGroup
		for i, item := range items {
			if itemMap, itemIsMap := item.(map[string]interface{}); itemIsMap {
				wg.Add(1)
				go func(i int, itemMap map[string]interface{}) {
					defer wg.Done()
					var processedItem map[string]interface{}
					if replaceItemObj {
						processedItem = replaceItem(itemMap, service)
					} else {
						processedItem = processItem(itemMap, service)
					}

					data["items"].([]interface{})[i] = processedItem

					resultChannel <- processedItem
				}(i, itemMap)
			}
		}

		go func() {
			wg.Wait()
			close(resultChannel)
		}()
	}

	for processedItem := range resultChannel {
		_ = append(processedItems, processedItem)
	}

	return data
}

func processItem(itemMap map[string]interface{}, service *youtube.Service) map[string]interface{} {
	snippet, snippetExists := itemMap["snippet"].(map[string]interface{})
	if !snippetExists {
		snippet = make(map[string]interface{})
		itemMap["snippet"] = snippet
	}

	channelId := snippet["channelId"]

	if channelIdStr, ok := channelId.(string); ok {
		channelThumbnails := getChannelThumbnails(service, channelIdStr)
		snippet["channelThumbnails"] = channelThumbnails
	}
	return itemMap
}

func replaceItem(itemMap map[string]interface{}, service *youtube.Service) map[string]interface{} {
	videoIdObj, videoIdExists := itemMap["id"].(map[string]interface{})
	if !videoIdExists {
		videoIdObj = make(map[string]interface{})
		itemMap["id"] = videoIdObj
	}

	videoId := videoIdObj["videoId"]

	if videoIdStr, ok := videoId.(string); ok {
		videoResponse := getVideo(service, videoIdStr)

		video := videoResponse.Items[0]

		item := make(map[string]interface{})
		item["kind"] = video.Kind
		item["etag"] = video.Etag
		item["id"] = video.Id
		item["contentDetails"] = video.ContentDetails
		item["snippet"] = video.Snippet
		item["statistics"] = video.Statistics
		item["topicDetails"] = video.TopicDetails

		itemMap = item
	}
	return itemMap
}
