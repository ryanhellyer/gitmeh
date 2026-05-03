//go:build !integration

package config

import (
	"testing"
)

func TestLoad_defaultChatWhenNoUserKey(t *testing.T) {
	t.Setenv("GITMEH_API_KEY", "")
	t.Setenv("OPENROUTER_API_KEY", "")
	t.Setenv("GITMEH_API_BASE", "")
	t.Setenv("GITMEH_MODEL", "")
	t.Setenv("OPENROUTER_MODEL", "")
	t.Setenv("GITMEH_PROMPT", "")

	got := Load()
	if got.Backend != BackendOpenAIChat {
		t.Fatalf("backend: got %v want chat", got.Backend)
	}
	if got.Chat.BaseURL != DefaultHostedChatBaseURL {
		t.Fatalf("BaseURL: got %q", got.Chat.BaseURL)
	}
	if got.Chat.APIKey != DefaultPublicAPIKey {
		t.Fatalf("APIKey: got %q want DefaultPublicAPIKey", got.Chat.APIKey)
	}
	if got.Chat.Model != DefaultHostedModel {
		t.Fatalf("Model: got %q", got.Chat.Model)
	}
	if got.Chat.Prompt == "" {
		t.Fatal("expected default prompt")
	}
}

func TestLoad_userAPIKeyOverridesDefaultPublic(t *testing.T) {
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

func TestLoad_customBaseWithDefaultPublicKey(t *testing.T) {
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
	if got.Chat.APIKey != DefaultPublicAPIKey {
		t.Fatalf("APIKey: got %q", got.Chat.APIKey)
	}
	if got.Chat.Model != DefaultHostedModel {
		t.Fatalf("Model: got %q", got.Chat.Model)
	}
}

func TestLoad_chatOpenRouterKey(t *testing.T) {
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

func TestLoad_fallbackModels(t *testing.T) {
	t.Setenv("GITMEH_API_KEY", "k")
	t.Setenv("GITMEH_FALLBACK_MODELS", " model-b,model-c , ")
	t.Setenv("GITMEH_API_BASE", "")
	t.Setenv("GITMEH_MODEL", "")
	t.Setenv("OPENROUTER_MODEL", "")
	t.Setenv("GITMEH_PROMPT", "")

	got := Load()
	if len(got.Chat.FallbackModels) != 2 {
		t.Fatalf("expected 2 fallback models, got %d: %v", len(got.Chat.FallbackModels), got.Chat.FallbackModels)
	}
	if got.Chat.FallbackModels[0] != "model-b" {
		t.Errorf("fallback[0] = %q", got.Chat.FallbackModels[0])
	}
	if got.Chat.FallbackModels[1] != "model-c" {
		t.Errorf("fallback[1] = %q", got.Chat.FallbackModels[1])
	}
}

func TestLoad_fallbackModelsEmptyEnv(t *testing.T) {
	t.Setenv("GITMEH_API_KEY", "k")
	t.Setenv("GITMEH_FALLBACK_MODELS", "")
	t.Setenv("GITMEH_API_BASE", "")
	t.Setenv("GITMEH_MODEL", "")

	got := Load()
	if len(got.Chat.FallbackModels) != 0 {
		t.Fatalf("expected 0 fallback models, got %d", len(got.Chat.FallbackModels))
	}
}

func TestLoad_chatOpenRouterModelEnv(t *testing.T) {
	t.Setenv("GITMEH_API_KEY", "x")
	t.Setenv("GITMEH_MODEL", "")
	t.Setenv("OPENROUTER_MODEL", "anthropic/claude-3-haiku")

	got := Load()
	if got.Chat.Model != "anthropic/claude-3-haiku" {
		t.Fatalf("Model: got %q", got.Chat.Model)
	}
}
