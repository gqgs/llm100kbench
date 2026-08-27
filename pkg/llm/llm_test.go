package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gqgs/llminvestbench/pkg/modelconfig"
	"github.com/stretchr/testify/require"
)

type fakeHTTPClient struct {
	responses []*http.Response
	requests  []*http.Request
	calls     int
}

func (f *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	f.requests = append(f.requests, req)
	response := f.responses[f.calls]
	f.calls++
	return response, nil
}

func TestProviderClientsRequestJSONMode(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		model    string
	}{
		{name: "mistral", provider: "mistral", model: "mistral-small-2603"},
		{name: "qwen", provider: "groq", model: "qwen/qwen3.8-27b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &fakeHTTPClient{responses: []*http.Response{
				response(http.StatusOK, `{"choices":[{"message":{"content":"{\"updates\":[],\"context\":[\"hold\"]}"}}]}`, nil),
			}}
			created, err := New(modelconfig.Model{Provider: tt.provider, Model: tt.model}, "key")
			require.NoError(t, err)
			client := created.(*openAICompatible)
			client.client = httpClient

			_, err = client.Generate(context.Background(), "return JSON")
			require.NoError(t, err)

			var body map[string]any
			require.NoError(t, json.NewDecoder(httpClient.requests[0].Body).Decode(&body))
			require.Equal(t, map[string]any{"type": "json_object"}, body["response_format"])
			if tt.name == "qwen" {
				require.Equal(t, "none", body["reasoning_effort"])
				require.Equal(t, "hidden", body["reasoning_format"])
			}
		})
	}
}

func TestOpenAICompatibleRetriesRateLimitUsingRetryAfter(t *testing.T) {
	client := &fakeHTTPClient{responses: []*http.Response{
		response(http.StatusTooManyRequests, "limited", http.Header{"Retry-After": []string{"2"}}),
		response(http.StatusOK, `{"choices":[{"message":{"content":"{\"updates\":[],\"context\":[\"hold\"]}"}}]}`, nil),
	}}
	var delays []time.Duration
	llm := &openAICompatible{
		endpoint: "https://example.invalid",
		model:    "model",
		apiKey:   "key",
		client:   client,
		wait: func(ctx context.Context, delay time.Duration) error {
			delays = append(delays, delay)
			return nil
		},
	}

	text, err := llm.Generate(context.Background(), "prompt")

	require.NoError(t, err)
	require.Contains(t, text, `"updates":[]`)
	require.Equal(t, 2, client.calls)
	require.Equal(t, []time.Duration{2 * time.Second}, delays)
}

func TestGeminiRetriesServerErrorUsingBackoff(t *testing.T) {
	client := &fakeHTTPClient{responses: []*http.Response{
		response(http.StatusServiceUnavailable, "unavailable", nil),
		response(http.StatusOK, `{"candidates":[{"content":{"parts":[{"text":"{\"updates\":[],\"context\":[\"hold\"]}"}]}}]}`, nil),
	}}
	var delays []time.Duration
	llm := &gemini{
		endpoint: "https://example.invalid",
		client:   client,
		wait: func(ctx context.Context, delay time.Duration) error {
			delays = append(delays, delay)
			return nil
		},
	}

	_, err := llm.Generate(context.Background(), "prompt")

	require.NoError(t, err)
	require.Equal(t, 2, client.calls)
	require.Equal(t, []time.Duration{time.Second}, delays)
}

func TestOpenAICompatibleStopsAfterTransientRetryLimit(t *testing.T) {
	client := &fakeHTTPClient{responses: []*http.Response{
		response(http.StatusTooManyRequests, "limited", nil),
		response(http.StatusTooManyRequests, "limited", nil),
		response(http.StatusTooManyRequests, "limited", nil),
	}}
	llm := &openAICompatible{
		endpoint: "https://example.invalid",
		model:    "model",
		client:   client,
		wait:     func(ctx context.Context, delay time.Duration) error { return nil },
	}

	_, err := llm.Generate(context.Background(), "prompt")

	require.Error(t, err)
	require.Contains(t, err.Error(), "429")
	require.Equal(t, maxRequestAttempts, client.calls)
}

func response(status int, body string, header http.Header) *http.Response {
	if header == nil {
		header = http.Header{}
	}
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
