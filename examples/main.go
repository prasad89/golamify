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
	}
	var wg sync.WaitGroup

	for i := 1; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			prompt := fmt.Sprintf("What is the square of %d?", i)
			resp, err := golamify.Generate(client, "llama3.2:1b", prompt)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(i,resp.Response)
		}(i)
	}

	wg.Wait()

	return
}
