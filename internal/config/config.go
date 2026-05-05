package config

import (
	"os"
	"strconv"
	"strings"
)

// Defaults for the built-in hosted OpenAI-compatible API (no user API key).
const (
	DefaultHostedChatBaseURL = "https://ai.hellyer.test/v1"
	DefaultPublicAPIKey      = "gitmeh-public-client" //nolint:gosec // public key for the default hosted endpoint
	DefaultHostedModel       = "gitmeh-hosted"
)

// Backend selects how gitmeh talks to the model service.
type Backend int

const (
	// BackendOpenAIChat uses an OpenAI-compatible /v1/chat/completions JSON API.
	BackendOpenAIChat Backend = iota
)

const DefaultMaxDiffBytes = 10_000

// OpenAIChat holds settings for [BackendOpenAIChat].
type OpenAIChat struct {
	BaseURL        string
	APIKey         string
	Model          string
	Prompt         string
	FallbackModels []string
	MaxDiffBytes   int // max staged diff size before refusing; 0 = no limit
}

// App is resolved configuration from the environment.
type App struct {
	Backend Backend
	Chat    OpenAIChat
}

const defaultOpenAIBase = "https://openrouter.ai/api/v1"
const defaultModel = "google/gemma-3-4b-it"

const defaultCommitPrompt = `Write a Git commit message (Conventional Commits format) for this diff. Reply with ONLY the commit message. No analysis, no explanation, no preamble. Start with a verb. No numbering. No bullet points.`

// Load reads environment variables.
//
// With GITMEH_API_KEY set: GITMEH_API_BASE defaults to OpenRouter, model to
// [defaultModel] unless GITMEH_MODEL is set.
//
// With no key set: [DefaultHostedChatBaseURL], [DefaultPublicAPIKey], and
// [DefaultHostedModel] unless GITMEH_API_BASE and/or GITMEH_MODEL override the
// URL or model.
//
// GITMEH_PROMPT optionally overrides the system instructions.
//
// OPENROUTER_API_KEY and OPENROUTER_MODEL are also accepted as legacy aliases
// for GITMEH_API_KEY and GITMEH_MODEL (from earlier versions that only
// supported OpenRouter).
func Load() App {
	userKey := strings.TrimSpace(os.Getenv("GITMEH_API_KEY"))
	if userKey == "" {
		userKey = strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY")) // legacy alias
	}

	prompt := strings.TrimSpace(os.Getenv("GITMEH_PROMPT"))
	if prompt == "" {
		prompt = defaultCommitPrompt
	}

	base := strings.TrimSpace(os.Getenv("GITMEH_API_BASE"))
	base = strings.TrimRight(base, "/")

	model := strings.TrimSpace(os.Getenv("GITMEH_MODEL"))
	if model == "" {
		model = strings.TrimSpace(os.Getenv("OPENROUTER_MODEL")) // legacy alias
	}

	fallbackRaw := strings.TrimSpace(os.Getenv("GITMEH_FALLBACK_MODELS"))
	var fallbackModels []string
	if fallbackRaw != "" {
		for _, p := range strings.Split(fallbackRaw, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				fallbackModels = append(fallbackModels, p)
			}
		}
	}

	maxDiff := DefaultMaxDiffBytes
	if raw := strings.TrimSpace(os.Getenv("GITMEH_MAX_DIFF_BYTES")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v >= 0 {
			maxDiff = v
		}
	}

	var apiKey string
	if userKey != "" {
		apiKey = userKey
		if base == "" {
			base = defaultOpenAIBase
		}
		if model == "" {
			model = defaultModel
		}
	} else {
		apiKey = DefaultPublicAPIKey
		if base == "" {
			base = DefaultHostedChatBaseURL
		}
		if model == "" {
			model = DefaultHostedModel
		}
	}

	return App{
		Backend: BackendOpenAIChat,
		Chat: OpenAIChat{
			BaseURL:        base,
			APIKey:         apiKey,
			Model:          model,
			Prompt:         prompt,
			FallbackModels: fallbackModels,
			MaxDiffBytes:   maxDiff,
		},
	}
}
