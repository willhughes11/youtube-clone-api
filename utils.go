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
		return nil // Handle error appropriately in your code
	}

	return channelResponse.Items[0].Snippet.Thumbnails
}

func processItemsConcurrently(jsonResponse []byte, service *youtube.Service) []map[string]interface{} {
	// Unmarshal the JSON into a Go data structure (e.g., a map)
	var data map[string]interface{}
	if err := json.Unmarshal(jsonResponse, &data); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return nil
	}

	// Create a channel to collect results from goroutines
	resultChannel := make(chan map[string]interface{})

	// Slice to store processed items
	processedItems := []map[string]interface{}{}

	// Loop through the data and launch a goroutine for each item
	if items, ok := data["items"].([]interface{}); ok {
		var wg sync.WaitGroup
		for _, item := range items {
			if itemMap, itemIsMap := item.(map[string]interface{}); itemIsMap {
				wg.Add(1)
				go func(itemMap map[string]interface{}) {
					defer wg.Done()
					processedItem := processItem(itemMap, service)
					resultChannel <- processedItem // Send the processed item to the channel
				}(itemMap)
			}
		}

		// Close the result channel when all goroutines are done
		go func() {
			wg.Wait()
			close(resultChannel)
		}()
	}

	// Collect results from goroutines
	for processedItem := range resultChannel {
		// Append processed items to the slice
		processedItems = append(processedItems, processedItem)
	}

	return processedItems
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

		// Modify or add a new key-value pair to "snippet"
		snippet["channelThumbnails"] = channelThumbnails
	}
	return itemMap
}
