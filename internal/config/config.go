package config

import (
	"os"
	"strings"
)

// DefaultPlainURL is the legacy keyless hosted endpoint: POST unified diff as
// text/plain; response body is the commit message. Used when GITMEH_LEGACY_PLAIN
// is enabled.
const DefaultPlainURL = "https://ai.hellyer.kiwi/gitmeh"

// HostedChatBaseURL is the default OpenAI-compatible API root for the built-in
// hosted service (no trailing slash). POST {HostedChatBaseURL}/chat/completions.
const HostedChatBaseURL = "https://ai.hellyer.kiwi/v1"

// HostedPublicBearer is the default Authorization bearer for the hosted chat
// endpoint. It is a weak client identifier, not a billing secret; the server
// should enforce limits by IP as today. See docs/hosted-api-migration-instructions.md.
const HostedPublicBearer = "gitmeh-public-client"

// HostedDefaultModel is sent in JSON for the hosted path; the server may ignore
// it and always use its local model.
const HostedDefaultModel = "gitmeh-hosted"

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
// Default (no user API key): OpenAI-compatible chat against [HostedChatBaseURL]
// with [HostedPublicBearer] unless GITMEH_HOSTED_TOKEN is set. Override base with
// GITMEH_API_BASE for staging.
//
// User API key (GITMEH_API_KEY or OPENROUTER_API_KEY): chat against
// GITMEH_API_BASE or OpenRouter; model from GITMEH_MODEL / OPENROUTER_MODEL.
//
// Legacy plain (GITMEH_LEGACY_PLAIN=true): POST text/plain to GITMEH_DEFAULT_URL
// or [DefaultPlainURL].
//
// GITMEH_PROMPT optionally overrides the system instructions (chat modes only).
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

	if userKey != "" {
		if base == "" {
			base = defaultOpenAIBase
		}
		if model == "" {
			model = defaultModel
		}
		return App{
			Backend:  BackendOpenAIChat,
			PlainURL: plain,
			Chat: OpenAIChat{
				BaseURL: base,
				APIKey:  userKey,
				Model:   model,
				Prompt:  prompt,
			},
		}
	}

	if base == "" {
		base = HostedChatBaseURL
	}
	hostedKey := strings.TrimSpace(os.Getenv("GITMEH_HOSTED_TOKEN"))
	if hostedKey == "" {
		hostedKey = HostedPublicBearer
	}
	if model == "" {
		model = HostedDefaultModel
	}
	return App{
		Backend:  BackendOpenAIChat,
		PlainURL: plain,
		Chat: OpenAIChat{
			BaseURL: base,
			APIKey:  hostedKey,
			Model:   model,
			Prompt:  prompt,
		},
	}
}
