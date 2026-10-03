package groq

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	endpoint = "https://api.groq.com/openai/v1/chat/completions"
	model    = "qwen/qwen3.8-27b"
	timeout  = 30 * time.Second

	maxTokens = 1024
)

var ErrRateLimited = errors.New("groq: rate limited")

type Client struct {
	key  string
	http *http.Client
}

func New(key string) *Client {
	return &Client{key: key, http: &http.Client{Timeout: timeout}}
}

func (c *Client) Configured() bool {
	return c.key != ""
}

type message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type part struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL string `json:"url"`
}

type request struct {
	Model          string            `json:"model"`
	Messages       []message         `json:"messages"`
	Temperature    float64           `json:"temperature"`
	MaxTokens      int               `json:"max_completion_tokens"`
	ResponseFormat map[string]string `json:"response_format"`
}

type response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Vision asks the model about an image and returns its answer, a JSON object.
func (c *Client) Vision(ctx context.Context, prompt string, image []byte) (string, error) {
	if !c.Configured() {
		return "", errors.New("groq: GROQ_KEY is not set")
	}

	dataURL := "data:" + http.DetectContentType(image) + ";base64," + base64.StdEncoding.EncodeToString(image)

	body, err := json.Marshal(request{
		Model: model,
		Messages: []message{
			{Role: "system", Content: prompt},
			{Role: "user", Content: []part{{Type: "image_url", ImageURL: &imageURL{URL: dataURL}}}},
		},
		MaxTokens:      maxTokens,
		ResponseFormat: map[string]string{"type": "json_object"},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("groq: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("groq: %w", err)
	}
	if res.StatusCode == http.StatusTooManyRequests {
		return "", ErrRateLimited
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq: status %d: %s", res.StatusCode, data)
	}

	var parsed response
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", fmt.Errorf("groq: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return "", errors.New("groq: empty response")
	}

	return parsed.Choices[0].Message.Content, nil
}
