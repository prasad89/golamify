# GoLamify

The **GoLamify** Go package provides an easy way to integrate Go projects with **Ollama**.

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

	resp, err := golamify.Generate(client, "llama3.2", "Why is the sky blue?")
	if err != nil {
		fmt.Println("Error generating response:", err)
		return
	}

	fmt.Println("Response:", resp.Response)
}
```

### 📂 More Examples

Explore additional examples in the `examples` directory to see how you can make the most of GoLamify.

## 👍 Contributing

Help us make GoLamify even better:

- Star this repo on GitHub! 🌟
- Submit issues and pull requests for improvements and bug fixes.
