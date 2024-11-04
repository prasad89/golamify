package golamify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GeneratePayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason"`
	Context            []int  `json:"context"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int64  `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int64  `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

func Generate(c *Client, model string, prompt string) (*GenerateResponse, error) {
	statusCode, err := ShowModel(model, c)
	if err != nil {
		return nil, fmt.Errorf("error showing model: %w", err)
	}

	if statusCode == http.StatusNotFound {
		pullStatus, err := PullModel(model, c)
		if err != nil {
			return nil, fmt.Errorf("failed to pull model: %w", err)
		}
		if pullStatus != http.StatusOK {
			return nil, fmt.Errorf("failed to pull model, received status: %d", pullStatus)
		}
	} else if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from ShowModel: %d", statusCode)
	}

	payload := GeneratePayload{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.OllamaHost+"/api/generate", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to generate endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("generate endpoint is not reachable, received status: %s, body: %s", resp.Status, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var generateResponse GenerateResponse
	if err := json.Unmarshal(respBody, &generateResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &generateResponse, nil
}
