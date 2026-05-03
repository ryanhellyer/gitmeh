package aiapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OpenAIChatParams configures a call to an OpenAI-compatible chat completions API.
type OpenAIChatParams struct {
	BaseURL      string // e.g. https://openrouter.ai/api/v1 (no trailing slash)
	APIKey       string
	Model        string
	SystemPrompt string // instructions; diff is sent in a separate user message
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
// text trimmed of outer whitespace. client and API key must be non-nil / non-empty.
func CommitMessageOpenAIChat(client *http.Client, p OpenAIChatParams, diff string) (string, error) {
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

	body := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: sys},
			{Role: "user", Content: "Unified diff:\n" + diff},
		},
		Temperature: 0.3,
		MaxTokens:   512,
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	endpoint := base + "/chat/completions"
	return withGeneratingCommitSpinner(func() (string, error) {
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(rawBody))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+key)
		if strings.Contains(strings.ToLower(base), "openrouter.ai") {
			req.Header.Set("HTTP-Referer", "https://github.com/ryanhellyer/gitmeh")
			req.Header.Set("X-Title", "gitmeh")
		}

		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

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
	})
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
