package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/prasad89/golamify/pkg/golamify"
)

func main() {
	config := golamify.Config{
		OllamaHost: "http://localhost:11434",
		Timeout:    300 * time.Second,
	}

	client, err := golamify.NewClient(&config)
	if err != nil {
		fmt.Print(err)
		return
	}

	var wg sync.WaitGroup

	for i := 1; i <= 1; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			prompt := fmt.Sprintf("What is the square of %d?", i)

			payload := golamify.GeneratePayload{
				Model:  "llama3.2:1b",
				Prompt: prompt,
				Stream: new(bool),
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
					}
				}
			}
		}(i)
	}

	wg.Wait()
}
