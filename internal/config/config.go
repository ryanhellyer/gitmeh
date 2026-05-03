package config

import (
	"os"
	"strings"
)

// DefaultPlainURL is the legacy keyless hosted endpoint: POST unified diff as
// text/plain; response body is the commit message. Used when GITMEH_LEGACY_PLAIN
// is enabled.
const DefaultPlainURL = "https://ai.hellyer.kiwi/gitmeh"

// Defaults for the built-in hosted OpenAI-compatible API (no user API key).
// Override with GITMEH_API_BASE / GITMEH_MODEL / GITMEH_API_KEY or OPENROUTER_API_KEY.
const (
	DefaultHostedChatBaseURL = "https://ai.hellyer.test/v1"
	DefaultPublicAPIKey      = "gitmeh-public-client"
	DefaultHostedModel       = "gitmeh-hosted"
)

// Backend selects how gitmeh talks to the model service.
type Backend int

const (
	// BackendPlain POSTs the diff to PlainURL as text/plain (no API key).
	BackendPlain Backend = iota
	// BackendOpenAIChat uses an OpenAI-compatible /v1/chat/completions JSON API.
	BackendOpenAIChat
)

// OpenAIChat holds settings for [BackendOpenAIChat].
type OpenAIChat struct {
	BaseURL string
	APIKey  string
	Model   string
	Prompt  string
}

// App is resolved configuration from the environment.
type App struct {
	Backend  Backend
	PlainURL string
	Chat     OpenAIChat
}

const defaultOpenAIBase = "https://openrouter.ai/api/v1"
const defaultModel = "google/gemma-3-4b-it"

const defaultCommitPrompt = `You write git commit messages.
Write a single concise commit message for the following unified diff (subject line only, or subject plus body if the change truly needs it).
Respond with only the commit message text. No markdown fences, no quotes, no preamble.`

func legacyPlainEnabled() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("GITMEH_LEGACY_PLAIN")))
	return v == "1" || v == "true" || v == "yes"
}

// Load reads environment variables.
//
// With GITMEH_API_KEY or OPENROUTER_API_KEY: GITMEH_API_BASE defaults to OpenRouter,
// model to [defaultModel] unless GITMEH_MODEL / OPENROUTER_MODEL is set.
//
// With neither key set: [DefaultHostedChatBaseURL], [DefaultPublicAPIKey], and
// [DefaultHostedModel] unless GITMEH_API_BASE and/or GITMEH_MODEL / OPENROUTER_MODEL
// override the URL or model.
//
// Legacy plain (GITMEH_LEGACY_PLAIN=true): POST text/plain to GITMEH_DEFAULT_URL
// or [DefaultPlainURL].
//
// GITMEH_PROMPT optionally overrides the system instructions (chat mode only).
func Load() App {
	plain := strings.TrimSpace(os.Getenv("GITMEH_DEFAULT_URL"))
	if plain == "" {
		plain = DefaultPlainURL
	}

	if legacyPlainEnabled() {
		return App{Backend: BackendPlain, PlainURL: plain}
	}

	userKey := strings.TrimSpace(os.Getenv("GITMEH_API_KEY"))
	if userKey == "" {
		userKey = strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY"))
	}

	prompt := strings.TrimSpace(os.Getenv("GITMEH_PROMPT"))
	if prompt == "" {
		prompt = defaultCommitPrompt
	}

	base := strings.TrimSpace(os.Getenv("GITMEH_API_BASE"))
	base = strings.TrimRight(base, "/")

	model := strings.TrimSpace(os.Getenv("GITMEH_MODEL"))
	if model == "" {
		model = strings.TrimSpace(os.Getenv("OPENROUTER_MODEL"))
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
		Backend:  BackendOpenAIChat,
		PlainURL: plain,
		Chat: OpenAIChat{
			BaseURL: base,
			APIKey:  apiKey,
			Model:   model,
			Prompt:  prompt,
		},
	}
}
