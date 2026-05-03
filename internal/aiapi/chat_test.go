//go:build !integration

package aiapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCommitMessageOpenAIChat_nilClient(t *testing.T) {
	t.Parallel()

	_, err := CommitMessageOpenAIChat(nil, OpenAIChatParams{
		BaseURL:      "http://x",
		APIKey:       "k",
		Model:        "m",
		SystemPrompt: "p",
	}, "diff")
	if err == nil || !strings.Contains(err.Error(), "nil") {
		t.Fatalf("got %v", err)
	}
}

func TestCommitMessageOpenAIChat_emptyKey(t *testing.T) {
	t.Parallel()

	_, err := CommitMessageOpenAIChat(DefaultHTTPClient(), OpenAIChatParams{
		BaseURL:      "http://x",
		APIKey:       "  ",
		Model:        "m",
		SystemPrompt: "p",
	}, "diff")
	if err == nil || !strings.Contains(err.Error(), "api key") {
		t.Fatalf("got %v", err)
	}
}

func TestCommitMessageOpenAIChat_success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method: %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("Authorization: %q", r.Header.Get("Authorization"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var req chatRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatal(err)
		}
		if req.Model != "test-model" {
			t.Errorf("model: %q", req.Model)
		}
		if len(req.Messages) != 2 {
			t.Fatalf("messages len: %d", len(req.Messages))
		}
		if req.Messages[0].Role != "system" || req.Messages[0].Content != "sys-here" {
			t.Errorf("system msg: %+v", req.Messages[0])
		}
		if req.Messages[1].Role != "user" || !strings.Contains(req.Messages[1].Content, "diff-here") {
			t.Errorf("user msg: %+v", req.Messages[1])
		}

		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"  fix typo  \n"}}]}`))
	}))
	defer srv.Close()

	got, err := CommitMessageOpenAIChat(srv.Client(), OpenAIChatParams{
		BaseURL:      srv.URL,
		APIKey:       "secret",
		Model:        "test-model",
		SystemPrompt: "sys-here",
	}, "diff-here")
	if err != nil {
		t.Fatal(err)
	}
	if got != "fix typo" {
		t.Fatalf("got %q want trimmed", got)
	}
}

func TestCommitMessageOpenAIChat_non2xxUsesErrorMessage(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
		_, _ = io.WriteString(w, `{"error":{"message":"insufficient_quota"}}`)
	}))
	defer srv.Close()

	_, err := CommitMessageOpenAIChat(srv.Client(), OpenAIChatParams{
		BaseURL:      srv.URL,
		APIKey:       "k",
		Model:        "m",
		SystemPrompt: "p",
	}, "d")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "402") {
		t.Errorf("want status in err: %v", err)
	}
	if !strings.Contains(err.Error(), "insufficient_quota") {
		t.Errorf("want API message: %v", err)
	}
}

func TestCommitMessageOpenAIChat_noChoices(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[]}`))
	}))
	defer srv.Close()

	_, err := CommitMessageOpenAIChat(srv.Client(), OpenAIChatParams{
		BaseURL:      srv.URL,
		APIKey:       "k",
		Model:        "m",
		SystemPrompt: "p",
	}, "d")
	if err == nil || !strings.Contains(err.Error(), "no choices") {
		t.Fatalf("got %v", err)
	}
}
