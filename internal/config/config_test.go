//go:build !integration

package config

import (
	"testing"
)

func TestLoad_builtinKeyWhenNoUserKey(t *testing.T) {
	prev := BuiltinAPIKey
	BuiltinAPIKey = "builtin-from-link-time"
	t.Cleanup(func() { BuiltinAPIKey = prev })

	t.Setenv("GITMEH_LEGACY_PLAIN", "")
	t.Setenv("GITMEH_API_KEY", "")
	t.Setenv("OPENROUTER_API_KEY", "")
	t.Setenv("GITMEH_API_BASE", "")
	t.Setenv("GITMEH_MODEL", "")
	t.Setenv("OPENROUTER_MODEL", "")
	t.Setenv("GITMEH_PROMPT", "")
	t.Setenv("GITMEH_DEFAULT_URL", "")

	got := Load()
	if got.Backend != BackendOpenAIChat {
		t.Fatalf("backend: got %v want chat", got.Backend)
	}
	if got.Chat.BaseURL != "https://openrouter.ai/api/v1" {
		t.Fatalf("BaseURL: got %q", got.Chat.BaseURL)
	}
	if got.Chat.APIKey != "builtin-from-link-time" {
		t.Fatalf("APIKey: got %q", got.Chat.APIKey)
	}
	if got.Chat.Model != "google/gemma-3-4b-it" {
		t.Fatalf("Model: got %q", got.Chat.Model)
	}
	if got.Chat.Prompt == "" {
		t.Fatal("expected default prompt")
	}
}

func TestLoad_userAPIKeyOverridesBuiltin(t *testing.T) {
	prev := BuiltinAPIKey
	BuiltinAPIKey = "builtin-from-link-time"
	t.Cleanup(func() { BuiltinAPIKey = prev })

	t.Setenv("GITMEH_LEGACY_PLAIN", "")
	t.Setenv("GITMEH_API_KEY", "user-key")
	t.Setenv("OPENROUTER_API_KEY", "")
	t.Setenv("GITMEH_API_BASE", "")
	t.Setenv("GITMEH_MODEL", "")
	t.Setenv("OPENROUTER_MODEL", "")

	got := Load()
	if got.Chat.APIKey != "user-key" {
		t.Fatalf("APIKey: got %q want user override", got.Chat.APIKey)
	}
}

func TestLoad_builtinKeyCustomBase(t *testing.T) {
	prev := BuiltinAPIKey
	BuiltinAPIKey = "builtin-from-link-time"
	t.Cleanup(func() { BuiltinAPIKey = prev })

	t.Setenv("GITMEH_LEGACY_PLAIN", "")
	t.Setenv("GITMEH_API_KEY", "")
	t.Setenv("OPENROUTER_API_KEY", "")
	t.Setenv("GITMEH_API_BASE", "https://staging.example/v1")
	t.Setenv("GITMEH_MODEL", "")
	t.Setenv("OPENROUTER_MODEL", "")
	t.Setenv("GITMEH_PROMPT", "")

	got := Load()
	if got.Backend != BackendOpenAIChat {
		t.Fatalf("backend")
	}
	if got.Chat.BaseURL != "https://staging.example/v1" {
		t.Fatalf("BaseURL: got %q", got.Chat.BaseURL)
	}
	if got.Chat.APIKey != "builtin-from-link-time" {
		t.Fatalf("APIKey")
	}
}

func TestLoad_legacyPlainDefaultURL(t *testing.T) {
	t.Setenv("GITMEH_LEGACY_PLAIN", "true")
	t.Setenv("GITMEH_API_KEY", "")
	t.Setenv("OPENROUTER_API_KEY", "")
	t.Setenv("GITMEH_DEFAULT_URL", "")

	got := Load()
	if got.Backend != BackendPlain {
		t.Fatalf("backend: got %v want plain", got.Backend)
	}
	if got.PlainURL != DefaultPlainURL {
		t.Fatalf("PlainURL: got %q", got.PlainURL)
	}
}

func TestLoad_legacyPlainCustomURL(t *testing.T) {
	t.Setenv("GITMEH_LEGACY_PLAIN", "1")
	t.Setenv("GITMEH_API_KEY", "")
	t.Setenv("OPENROUTER_API_KEY", "")
	t.Setenv("GITMEH_DEFAULT_URL", "https://example.com/git")

	got := Load()
	if got.Backend != BackendPlain {
		t.Fatalf("backend: got %v want plain", got.Backend)
	}
	if got.PlainURL != "https://example.com/git" {
		t.Fatalf("PlainURL: got %q", got.PlainURL)
	}
}

func TestLoad_chatOpenRouterKey(t *testing.T) {
	t.Setenv("GITMEH_LEGACY_PLAIN", "")
	t.Setenv("GITMEH_API_KEY", "")
	t.Setenv("OPENROUTER_API_KEY", "sk-test")
	t.Setenv("GITMEH_API_BASE", "")
	t.Setenv("GITMEH_MODEL", "")
	t.Setenv("OPENROUTER_MODEL", "")
	t.Setenv("GITMEH_PROMPT", "")

	got := Load()
	if got.Backend != BackendOpenAIChat {
		t.Fatalf("backend: got %v want chat", got.Backend)
	}
	if got.Chat.APIKey != "sk-test" {
		t.Fatalf("APIKey")
	}
	if got.Chat.BaseURL != "https://openrouter.ai/api/v1" {
		t.Fatalf("BaseURL: got %q", got.Chat.BaseURL)
	}
	if got.Chat.Model != "google/gemma-3-4b-it" {
		t.Fatalf("Model: got %q", got.Chat.Model)
	}
}

func TestLoad_chatGITMEHKeyOverridesBase(t *testing.T) {
	t.Setenv("GITMEH_LEGACY_PLAIN", "")
	t.Setenv("GITMEH_API_KEY", "k")
	t.Setenv("OPENROUTER_API_KEY", "")
	t.Setenv("GITMEH_API_BASE", "https://api.openai.com/v1")
	t.Setenv("GITMEH_MODEL", "gpt-4o-mini")
	t.Setenv("GITMEH_PROMPT", "Be brief.")

	got := Load()
	if got.Backend != BackendOpenAIChat {
		t.Fatalf("backend")
	}
	if got.Chat.APIKey != "k" {
		t.Fatalf("APIKey")
	}
	if got.Chat.BaseURL != "https://api.openai.com/v1" {
		t.Fatalf("BaseURL: got %q", got.Chat.BaseURL)
	}
	if got.Chat.Model != "gpt-4o-mini" {
		t.Fatalf("Model")
	}
	if got.Chat.Prompt != "Be brief." {
		t.Fatalf("Prompt: got %q", got.Chat.Prompt)
	}
}

func TestLoad_chatGITMEHKeyPreferredOverOpenRouter(t *testing.T) {
	t.Setenv("GITMEH_LEGACY_PLAIN", "")
	t.Setenv("GITMEH_API_KEY", "primary")
	t.Setenv("OPENROUTER_API_KEY", "secondary")
	t.Setenv("GITMEH_API_BASE", "")
	t.Setenv("GITMEH_MODEL", "")
	t.Setenv("OPENROUTER_MODEL", "")

	got := Load()
	if got.Chat.APIKey != "primary" {
		t.Fatalf("want GITMEH_API_KEY to win, got %q", got.Chat.APIKey)
	}
}

func TestLoad_chatOpenRouterModelEnv(t *testing.T) {
	t.Setenv("GITMEH_LEGACY_PLAIN", "")
	t.Setenv("GITMEH_API_KEY", "x")
	t.Setenv("GITMEH_MODEL", "")
	t.Setenv("OPENROUTER_MODEL", "anthropic/claude-3-haiku")

	got := Load()
	if got.Chat.Model != "anthropic/claude-3-haiku" {
		t.Fatalf("Model: got %q", got.Chat.Model)
	}
}
