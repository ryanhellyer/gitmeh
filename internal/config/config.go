package config

import (
	"os"
	"strings"
)

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
	// BackendOpenAIChat uses an OpenAI-compatible /v1/chat/completions JSON API.
	BackendOpenAIChat Backend = iota
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
	Backend Backend
	Chat    OpenAIChat
}

const defaultOpenAIBase = "https://openrouter.ai/api/v1"
const defaultModel = "google/gemma-3-4b-it"

const defaultCommitPrompt = `You write git commit messages.
Write a single concise commit message for the following unified diff (subject line only, or subject plus body if the change truly needs it).
Respond with only the commit message text. No markdown fences, no quotes, no preamble.`

// Load reads environment variables.
//
// With GITMEH_API_KEY or OPENROUTER_API_KEY: GITMEH_API_BASE defaults to OpenRouter,
// model to [defaultModel] unless GITMEH_MODEL / OPENROUTER_MODEL is set.
//
// With neither key set: [DefaultHostedChatBaseURL], [DefaultPublicAPIKey], and
// [DefaultHostedModel] unless GITMEH_API_BASE and/or GITMEH_MODEL / OPENROUTER_MODEL
// override the URL or model.
//
// GITMEH_PROMPT optionally overrides the system instructions.
func Load() App {
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
		Backend: BackendOpenAIChat,
		Chat: OpenAIChat{
			BaseURL: base,
			APIKey:  apiKey,
			Model:   model,
			Prompt:  prompt,
		},
	}
}
