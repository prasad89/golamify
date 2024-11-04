package golamify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PullResponse struct {
	Status string `json:"status"`
	Digest string `json:"digest"`
	Total  string `json:"total"`
	Error  string `josn:"error`
}

func ShowModel(model string, c *Client) (int, error) {
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

func PullModel(model string, c *Client) (int, error) {
	payload := map[string]string{"name": model}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to encode request payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.OllamaHost+"/api/pull", bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to pull endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("pull endpoint is not reachable, received status: %s, body: %s", resp.Status, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var pullResponse PullResponse
	if err := json.Unmarshal(respBody, &pullResponse); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	return resp.StatusCode, nil
}
