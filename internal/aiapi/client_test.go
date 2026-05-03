//go:build !integration

package aiapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPClientForChatBase_insecureTransportOnlyForDevHost(t *testing.T) {
	t.Parallel()

	dev := HTTPClientForChatBase("https://ai.hellyer.test/v1")
	if dev.Transport == nil {
		t.Fatal("expected custom transport for ai.hellyer.test")
	}
	prod := HTTPClientForChatBase("https://openrouter.ai/api/v1")
	if prod.Transport != nil {
		t.Fatal("expected default transport for non-dev host")
	}
}

func TestCommitMessage_nilClient(t *testing.T) {
	t.Parallel()

	_, err := CommitMessage(nil, "https://example.com/api", "diff")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Fatalf("error: %v", err)
	}
}

func TestCommitMessage_success_trimsBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "  hello from api  \n\t")
	}))
	defer srv.Close()

	got, err := CommitMessage(srv.Client(), srv.URL, "diff")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello from api" {
		t.Fatalf("got %q want trimmed message", got)
	}
}

func TestCommitMessage_postsExpectedRequest(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method: %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "text/plain; charset=UTF-8" {
			t.Errorf("Content-Type: %q", ct)
		}
		if ac := r.Header.Get("Accept"); ac != "text/plain" {
			t.Errorf("Accept: %q", ac)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != "unified-diff-bytes" {
			t.Errorf("body: %q", body)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	_, err := CommitMessage(srv.Client(), srv.URL, "unified-diff-bytes")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommitMessage_non2xxIncludesStatusAndBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "OpenRouter sorry", http.StatusBadGateway)
	}))
	defer srv.Close()

	_, err := CommitMessage(srv.Client(), srv.URL, "x")
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "502") {
		t.Errorf("expected status in error, got: %s", msg)
	}
	if !strings.Contains(msg, "OpenRouter sorry") {
		t.Errorf("expected raw body substring in error, got: %s", msg)
	}
	if !strings.Contains(msg, "| raw body:") {
		t.Errorf("expected formatter prefix in error, got: %s", msg)
	}
}
