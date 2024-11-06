# GoLamify

The **GoLamify** Go package provides an easy way to integrate Go projects with **Ollama**.

## ✨ Features

1. **Generate Responses from Ollama Models** – Easily generate responses using a variety of Ollama models.
2. **Default Response Streaming** – Real-time response streaming for immediate output.
3. **Full Parameter Support** – Customize model behavior with full API parameter support.
4. **No Model Pulling Needed** – Access models without manual pre-pulling.
5. **Clear Error Handling** – Simple, concise error handling for easy debugging.
6. **More** – Comming soon.

## 🚀 Getting Started

### Installation

To get started with GoLamify, add the following import to your code, and use Go’s module support to automatically fetch dependencies:

```go
import "github.com/prasad89/golamify/pkg/golamify"
```

Alternatively, install it using:

```bash
go get -u github.com/prasad89/golamify
```

### 🏃 Running GoLamify

Here's a simple example to get a GoLamify application up and running:

```go
package main

import (
	"fmt"

	"github.com/prasad89/golamify/pkg/golamify"
)

func main() {
	client, err := golamify.NewClient(nil)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}

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
```

### 📂 More Examples

Explore additional examples in the `examples` directory to see how you can make the most of GoLamify.

## 👍 Contributing

Help us make GoLamify even better:

- Star this repo on GitHub! 🌟
- Submit issues and pull requests for improvements and bug fixes.
