package aiapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// maxRetriesPerModel is the number of attempts per model for transient errors.
const maxRetriesPerModel = 3

// OpenAIChatParams configures a call to an OpenAI-compatible chat completions API.
type OpenAIChatParams struct {
	BaseURL        string   // e.g. https://openrouter.ai/api/v1 (no trailing slash)
	APIKey         string
	Model          string
	SystemPrompt   string   // instructions; diff is sent in a separate user message
	FallbackModels []string // models to try if primary fails with retryable error
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type apiErrorBody struct {
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// CommitMessageOpenAIChat POSTs to {BaseURL}/chat/completions with the unified diff
// in the user message and SystemPrompt as the system message. Returns assistant
// text trimmed of outer whitespace.
//
// If FallbackModels are provided and the primary model fails with a retryable
// error (network errors, 5xx, 429, context-length exceeded), fallback models
// are tried in order. Each model is retried up to [maxRetriesPerModel] times
// with exponential backoff for transient errors before moving to the next
// fallback.
func CommitMessageOpenAIChat(ctx context.Context, client *http.Client, p OpenAIChatParams, diff string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("http client is nil")
	}
	key := strings.TrimSpace(p.APIKey)
	if key == "" {
		return "", fmt.Errorf("api key is empty")
	}
	base := strings.TrimRight(strings.TrimSpace(p.BaseURL), "/")
	if base == "" {
		return "", fmt.Errorf("api base url is empty")
	}
	model := strings.TrimSpace(p.Model)
	if model == "" {
		return "", fmt.Errorf("model is empty")
	}
	sys := strings.TrimSpace(p.SystemPrompt)
	if sys == "" {
		return "", fmt.Errorf("system prompt is empty")
	}

	models := buildModelList(model, p.FallbackModels)

		return withGeneratingCommitSpinner(func() (string, error) {
		var lastErr error
		for i, m := range models {
			result, err := tryModelWithRetry(ctx, client, base, key, m, sys, diff)
			if err == nil {
				return result, nil
			}
			lastErr = err
			if i < len(models)-1 {
				fmt.Fprintf(os.Stderr, "\n  → trying fallback model %q ...\n", models[i+1])
			}
		}
		return "", &AllModelsFailedError{
			Models: models,
			Cause:  lastErr,
		}
	})
}

func buildModelList(primary string, fallbacks []string) []string {
	models := make([]string, 0, 1+len(fallbacks))
	models = append(models, primary)
	for _, m := range fallbacks {
		m = strings.TrimSpace(m)
		if m != "" && m != primary {
			models = append(models, m)
		}
	}
	return models
}

func tryModelWithRetry(ctx context.Context, client *http.Client, baseURL, apiKey, model, systemPrompt, diff string) (string, error) {
	var lastErr error
	for attempt := 0; attempt < maxRetriesPerModel; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<(attempt-1)) * time.Second
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return "", ctx.Err()
			case <-timer.C:
			}
		}
		result, err := doChatRequest(ctx, client, baseURL, apiKey, model, systemPrompt, diff)
		if err == nil {
			return result, nil
		}
		lastErr = err

		if isContextLengthError(err) {
			fmt.Fprintf(os.Stderr, "\n  %s: context length exceeded\n", model)
			return "", err
		}
		if !isRetryable(err) {
			fmt.Fprintf(os.Stderr, "\n  %s: %v\n", model, err)
			return "", err
		}
		fmt.Fprintf(os.Stderr, "\n  %s attempt %d/%d: %v\n", model, attempt+1, maxRetriesPerModel, err)
	}
	fmt.Fprintf(os.Stderr, "\n  %s failed after %d attempts\n", model, maxRetriesPerModel)
	return "", lastErr
}

func doChatRequest(ctx context.Context, client *http.Client, baseURL, apiKey, model, systemPrompt, diff string) (string, error) {
	body := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: "Unified diff:\n" + diff},
		},
		Temperature: 0.3,
		MaxTokens:   4096,
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	endpoint := baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(rawBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	if strings.Contains(strings.ToLower(baseURL), "openrouter.ai") {
		req.Header.Set("HTTP-Referer", "https://github.com/ryanhellyer/gitmeh")
		req.Header.Set("X-Title", "gitmeh")
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	raw := string(respBytes)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := summarizeChatAPIError(respBytes, raw)
		return "", fmt.Errorf("%s | %s", resp.Status, msg)
	}

	var parsed chatResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("decode response: %w (body: %q)", err, truncateForErr(raw))
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("no choices in response: %q", truncateForErr(raw))
	}
	out := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if out == "" {
		return "", fmt.Errorf("empty assistant content: %q", truncateForErr(raw))
	}
	return out, nil
}

func isRetryable(err error) bool {
	s := err.Error()
	if strings.Contains(s, "timeout") ||
		strings.Contains(s, "connection refused") ||
		strings.Contains(s, "no such host") ||
		strings.Contains(s, "connection reset") ||
		strings.Contains(s, "TLS handshake") {
		return true
	}
	if strings.Contains(s, "429") || strings.Contains(s, "500") || strings.Contains(s, "502") ||
		strings.Contains(s, "503") || strings.Contains(s, "504") {
		return true
	}
	if strings.Contains(s, "Provider returned error") {
		return true
	}
	return false
}

func isContextLengthError(err error) bool {
	s := err.Error()
	return strings.Contains(s, "maximum context length") ||
		strings.Contains(s, "context length") ||
		strings.Contains(s, "too many tokens")
}

// AllModelsFailedError is returned when every model (primary + fallbacks) fails.
type AllModelsFailedError struct {
	Models []string
	Cause  error
}

func (e *AllModelsFailedError) Error() string {
	return fmt.Sprintf("all %d models failed: %v", len(e.Models), e.Cause)
}

func (e *AllModelsFailedError) Unwrap() error {
	return e.Cause
}

func summarizeChatAPIError(respBytes []byte, raw string) string {
	var eb apiErrorBody
	if json.Unmarshal(respBytes, &eb) == nil && eb.Error != nil && strings.TrimSpace(eb.Error.Message) != "" {
		return eb.Error.Message
	}
	return fmt.Sprintf("raw body: %q", truncateForErr(raw))
}

func truncateForErr(s string) string {
	const max = 800
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
