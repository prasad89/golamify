package main

import (
	"fmt"
	"time"

	"github.com/prasad89/golamify/pkg/golamify"
)

func main() {
	// Optional config; pass nil for defaults
	config := golamify.Config{
		OllamaHost: "http://localhost:11434",
		Timeout:    30 * time.Second,
	}

	client, err := golamify.NewClient(&config)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}

	// Required parameters for Chat; others optional
	payload := golamify.ChatPayload{
		Model: "llama3.2:1b",
		Messages: []golamify.Message{
			{
				Role:    "user",
				Content: "What we discussed earlier?",
			},
		},
	}

	responseChannel, errorChannel := golamify.Chat(client, &payload)

	for {
		select {
		case response, ok := <-responseChannel:
			if !ok {
				return
			}

			// Check to safely extract "message" and "content"
			if messMap, exists := response["message"].(map[string]interface{}); exists {
				if content, exists := messMap["content"].(string); exists {
					fmt.Print(content)
				} else {
					fmt.Println("Content key missing or not a string")
				}
			} else {
				fmt.Println("Message key missing or not a valid map")
			}

		case err, ok := <-errorChannel:
			if ok && err != nil {
				fmt.Println("Error:", err)
			} else if !ok {
				return
			}
		}
	}
}
