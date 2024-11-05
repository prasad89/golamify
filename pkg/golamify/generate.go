package golamify

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GeneratePayload struct {
	Model     string   `json:"model" validate:"required"`
	Prompt    string   `json:"prompt" validate:"required"`
	Suffix    string   `json:"suffix,omitempty"`
	Images    []string `json:"images,omitempty"`
	System    string   `json:"system,omitempty"`
	Template  string   `json:"template,omitempty"`
	Context   string   `json:"context,omitempty"`
	Stream    *bool    `json:"stream,omitempty"`
	Raw       *bool    `json:"raw,omitempty"`
	KeepAlive string   `json:"keep_alive,omitempty"`
}

func Generate(c *Client, payload *GeneratePayload) (<-chan map[string]interface{}, <-chan error) {
	responseChannel := make(chan map[string]interface{})
	errorChannel := make(chan error, 1)

	go func() {
		defer close(responseChannel)
		defer close(errorChannel)

		statusCode, err := ShowModel(payload.Model, c)
		if err != nil {
			errorChannel <- fmt.Errorf("error showing model: %w", err)
			return
		}

		if statusCode == http.StatusNotFound {
			err := PullModel(payload.Model, c)
			if err != nil {
				errorChannel <- fmt.Errorf("failed to pull model: %w", err)
				return
			}
		}

		body, err := json.Marshal(payload)
		if err != nil {
			errorChannel <- fmt.Errorf("failed to encode request payload: %w", err)
			return
		}

		req, err := http.NewRequest("POST", c.config.OllamaHost+"/api/generate", bytes.NewBuffer(body))
		if err != nil {
			errorChannel <- fmt.Errorf("failed to create request: %w", err)
			return
		}
		req.Header.Set("User-Agent", c.userAgent)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			errorChannel <- fmt.Errorf("failed to connect to generate endpoint: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			errorChannel <- fmt.Errorf("generate endpoint is not reachable, received status: %s, body: %s", resp.Status, string(respBody))
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var generated map[string]interface{}
			err := json.Unmarshal(scanner.Bytes(), &generated)
			if err != nil {
				errorChannel <- fmt.Errorf("error parsing JSON: %w", err)
				continue
			}
			responseChannel <- generated
			if done, exists := generated["done"].(bool); exists && done {
				break
			}
		}

		if err := scanner.Err(); err != nil {
			errorChannel <- fmt.Errorf("error reading response body: %w", err)
		}
	}()

	return responseChannel, errorChannel
}
