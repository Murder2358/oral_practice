package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

type Client struct {
	apiKey      string
	baseURL     string
	model       string
	maxTokens   int
	temperature float64
	httpClient  *http.Client
}

func NewClient() *Client {
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	if apiKey == "" {
		panic("DASHSCOPE_API_KEY is not set in .env file")
	}

	return &Client{
		apiKey:      apiKey,
		baseURL:     "https://dashscope.aliyuncs.com/compatible-mode/v1",
		model:       "qwen-plus",
		maxTokens:   2048,
		temperature: 0.7,
		httpClient:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) WithModel(model string) *Client {
	c.model = model
	return c
}

func (c *Client) WithTemperature(temp float64) *Client {
	c.temperature = temp
	return c
}

func (c *Client) Chat(messages []Message) (string, error) {
	reqBody := chatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   c.maxTokens,
		Temperature: c.temperature,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return chatResp.Choices[0].Message.Content, nil
}
