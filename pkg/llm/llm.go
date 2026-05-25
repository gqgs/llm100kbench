package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gqgs/llminvestbench/pkg/modelconfig"
)

type Client interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type waitFunc func(context.Context, time.Duration) error

const (
	maxRequestAttempts = 3
	maxRetryDelay      = 30 * time.Second
)

func New(model modelconfig.Model, apiKey string) (Client, error) {
	switch model.Provider {
	case "github":
		return NewOpenAICompatible("https://models.github.ai/inference/chat/completions", model.Model, apiKey), nil
	case "groq":
		return NewOpenAICompatible("https://api.groq.com/openai/v1/chat/completions", model.Model, apiKey), nil
	case "mistral":
		return NewOpenAICompatible("https://api.mistral.ai/v1/chat/completions", model.Model, apiKey), nil
	case "gemini":
		return NewGemini(model.Model, apiKey), nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", model.Provider)
	}
}

type openAICompatible struct {
	endpoint string
	model    string
	apiKey   string
	client   httpClient
	wait     waitFunc
}

func NewOpenAICompatible(endpoint, model, apiKey string) Client {
	return &openAICompatible{
		endpoint: endpoint,
		model:    model,
		apiKey:   apiKey,
		client:   http.DefaultClient,
		wait:     waitContext,
	}
}

func (c *openAICompatible) Generate(ctx context.Context, prompt string) (string, error) {
	body := map[string]any{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": "Return only valid JSON. Do not wrap the answer in markdown."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(encoded))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := doWithRetry(ctx, c.client, req, c.wait)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("llm request failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if len(parsed.Choices) == 0 || strings.TrimSpace(parsed.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("llm response did not include message content")
	}
	return parsed.Choices[0].Message.Content, nil
}

type gemini struct {
	endpoint string
	apiKey   string
	client   httpClient
	wait     waitFunc
}

func NewGemini(model, apiKey string) Client {
	return &gemini{
		endpoint: "https://generativelanguage.googleapis.com/v1beta/models/" + model + ":generateContent?key=" + apiKey,
		apiKey:   apiKey,
		client:   http.DefaultClient,
		wait:     waitContext,
	}
}

func (c *gemini) Generate(ctx context.Context, prompt string) (string, error) {
	body := map[string]any{
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt + "\n\nReturn only valid JSON. Do not wrap the answer in markdown."},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature": 0.2,
		},
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(encoded))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := doWithRetry(ctx, c.client, req, c.wait)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("llm request failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini response did not include text")
	}
	text := strings.TrimSpace(parsed.Candidates[0].Content.Parts[0].Text)
	if text == "" {
		return "", fmt.Errorf("gemini response text was empty")
	}
	return text, nil
}

func doWithRetry(ctx context.Context, client httpClient, req *http.Request, wait waitFunc) (*http.Response, error) {
	if wait == nil {
		wait = waitContext
	}
	for attempt := 0; attempt < maxRequestAttempts; attempt++ {
		attemptReq := req
		if attempt > 0 {
			attemptReq = req.Clone(ctx)
			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return nil, err
				}
				attemptReq.Body = body
			}
		}

		resp, err := client.Do(attemptReq)
		if err != nil {
			return nil, err
		}
		if !isTransientStatus(resp.StatusCode) || attempt == maxRequestAttempts-1 {
			return resp, nil
		}

		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		_ = resp.Body.Close()
		if err := wait(ctx, retryDelay(resp, attempt)); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("llm request failed after retries")
}

func isTransientStatus(status int) bool {
	return status == http.StatusTooManyRequests || status >= 500 && status <= 599
}

func retryDelay(resp *http.Response, attempt int) time.Duration {
	if value := strings.TrimSpace(resp.Header.Get("Retry-After")); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil && seconds >= 0 {
			return boundedRetryDelay(time.Duration(seconds) * time.Second)
		}
		if at, err := http.ParseTime(value); err == nil {
			if delay := time.Until(at); delay > 0 {
				return boundedRetryDelay(delay)
			}
		}
	}
	return boundedRetryDelay(time.Duration(1<<attempt) * time.Second)
}

func boundedRetryDelay(delay time.Duration) time.Duration {
	if delay > maxRetryDelay {
		return maxRetryDelay
	}
	return delay
}

func waitContext(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
