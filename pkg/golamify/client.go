package golamify

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

const (
	defaultOllamaHost = "http://localhost:11434"
	defaultTimeout    = 30 * time.Second
)

var userAgent = "golamify-client/1.0 (" + runtime.GOARCH + " " + runtime.GOOS + ") Go/" + runtime.Version()

type Config struct {
	OllamaHost string
	Timeout    time.Duration
}

type Client struct {
	httpClient *http.Client
	config     Config
	userAgent  string
}

func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = &Config{}
	}

	if config.OllamaHost == "" {
		config.OllamaHost = defaultOllamaHost
	}
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	httpClient := &http.Client{Timeout: config.Timeout}

	client := &Client{
		httpClient: httpClient,
		config:     *config,
		userAgent:  userAgent,
	}

	if err := client.PingOllama(); err != nil {
		return nil, fmt.Errorf("unable to initialize client: %w", err)
	}

	return client, nil
}

func (c *Client) PingOllama() error {
	req, err := http.NewRequest("GET", c.config.OllamaHost, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to OllamaHost: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OllamaHost %s is not reachable, received status: %s", c.config.OllamaHost, resp.Status)
	}

	return nil
}
