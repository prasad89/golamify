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

	// Required parameters for Generate; others optional
	payload := golamify.GeneratePayload{
		Model:  "llama3.2:1b",
		Prompt: "Why is the sky blue?",
	}

	responseChannel, errorChannel := golamify.Generate(client, &payload)

	for {
		select {
		case response, ok := <-responseChannel:
			if !ok {
				return
			}
			fmt.Print(response["response"])

		case err, ok := <-errorChannel:
			if ok && err != nil {
				fmt.Println("Error:", err)
			} else if !ok {
				return
			}
		}
	}
}
