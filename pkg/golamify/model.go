package golamify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func ShowModel(c *Client, model string) (int, error) {
	payload := map[string]string{"name": model}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to encode request payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.OllamaHost+"/api/show", bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to show endpoint: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func PullModel(c *Client, model string) error {
	payload := map[string]string{"name": model}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode request payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.OllamaHost+"/api/pull", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to pull endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response body: %w", err)
		}
		return fmt.Errorf("pull endpoint is not reachable, received status: %s, body: %s", resp.Status, string(respBody))
	}

	fmt.Printf("Pulling model: %s", model)

	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var message struct {
			Error string `json:"error,omitempty"`
		}
		if err := decoder.Decode(&message); err != nil {
			return fmt.Errorf("failed to decode JSON message: %w", err)
		}

		if message.Error != "" {
			return fmt.Errorf("'%s': %s", model, message.Error)
		}

		fmt.Print("...")
		time.Sleep(5 * time.Second)
	}

	return nil
}
