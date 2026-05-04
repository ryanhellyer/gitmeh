# gitmeh ¯\_(ツ)_/¯

**AI-powered git commits for the terminally lazy.**

Stages everything (`git add --all`), AI-guesses a commit message, then shovels it to the cloud. Designed for personal garbage repos where quality does not matter and the only thing you care about is closing the laptop as fast as humanly possible.

> **⚠️** Not recommended for professional team projects. Using this at work is a great way to get a stern talking-to from your engineering manager.

### Why use this?

Because writing thoughtful commit messages for your 14th unfinished side project is a waste of your precious nap time.

* **Nuclear Staging:** It runs `git add --all` without asking. It stages your unfinished thoughts, your secrets, and that one large `test.mp4` you forgot was there.
* **AI Guesswork:** Generates a commit message via an OpenAI-compatible chat API, with retry logic and configurable fallback models.
* **Automatic Pushing:** Shovels your changes directly to the cloud so you can stop looking at the terminal.

### Default API service

If you don't set `GITMEH_API_KEY`, gitmeh uses a **free hosted API** at `https://ai.hellyer.test/`, run by the author (Ryan Hellyer). The backend automatically selects whichever AI model is working best and cheapest at the time, so models will vary between requests without warning.

**Your staged diff (code) is sent to this server** and then forwarded to whichever model the backend picks. If you are not comfortable sharing your code with me (Ryan) or with the random third-party model I route it through, **do not use the default service**. Set `GITMEH_API_BASE`, `GITMEH_API_KEY`, `GITMEH_MODEL` etc. to use your own AI provider instead.

I have zero interest in your code and no intention of looking at it, but it will be processed through my server and the model provider's servers.

## Quick Start

```bash
# 1. Install
make build && cp git-meh ~/.local/bin/           # from the repo root (requires Go)
# Or: ./install.sh                              # uses a prebuilt binary

# 2. Set up an API key (OpenCode Zen recommended)
export GITMEH_API_BASE='https://opencode.ai/zen/v1'
export GITMEH_API_KEY='your_zen_key'

# 3. Run
git meh
```

Git discovers the binary as a subcommand — works in any repository.

## Configuration

| Env var | Description | Default |
|---|---|---|
| `GITMEH_API_BASE` | API base URL | `https://ai.hellyer.test/v1` (built-in) |
| `GITMEH_API_KEY` | API key | built-in public key |
| `GITMEH_MODEL` | Model name | `gitmeh-hosted` or `google/gemma-3-4b-it` |
| `GITMEH_PROMPT` | System prompt for the model | Conventional Commits prompt |
| `GITMEH_FALLBACK_MODELS` | Comma-separated models to try if the primary fails | — |
| `GITMEH_MAX_DIFF_BYTES` | Per-file diff truncation limit (0 = no limit) | `10000` (10 KB) |

**Auth priority**: `GITMEH_API_KEY` > built-in public key.

**Fallback models**: If the primary model fails (timeout, 5xx, context-length exceeded), gitmeh retries up to 3 times with exponential backoff, then tries each fallback model in order. A 401 or other client error skips retries immediately.

**Diff truncation**: When the staged diff exceeds `GITMEH_MAX_DIFF_BYTES`, gitmeh keeps all file headers and proportionally trims hunk content per file. Truncated sections are marked with `# hunk truncated`.

## Developer Guide

### Prerequisites

- Go (see `go.mod` for version)
- `golangci-lint` and `govulncheck` for linting (install via `go install`)

### Commands

```bash
make build       # build native binary
make test        # run unit tests
make lint        # run golangci-lint + govulncheck
make cross       # cross-compile for Linux/macOS, amd64/arm64
make clean       # remove built binaries
make all         # lint + test + cross-compile

go test -tags=integration ./... -count=1   # integration tests (require git)
```

### Project structure

```
main.go              — entry point, CLI orchestration, user review prompt
internal/
  aiapi/             — AI API communication (chat, HTTP client, spinner)
  config/            — env var parsing
  git/               — git command wrappers (add, diff, commit, push)
  version/           — version string
```

### Architecture notes

- The API call wraps a spinner goroutine for terminal feedback. Ctrl+C cancels the HTTP context, which immediately aborts the request and cleans up the terminal.
- Model retries use exponential backoff (1s, 2s, 4s). Context-length errors and non-retryable status codes skip retries and advance to the next fallback model.
- Diff truncation splits the unified diff at `diff --git` boundaries, preserves all file headers, and allocates the remaining byte budget proportionally by hunk size.

## Changelog

- **3.x:** Retry and fallback models, graceful Ctrl+C, diff truncation, CI linting/security scanning, Dependabot, Makefile, support for OpenAI compatible APIs
- **3.0:** Rewrite in Go; run via `git meh`
- **2.x:** OpenRouter and plain-text API versions
- **1.0:** Initial Google Gemini implementation

## Author

**Ryan Hellyer** [ryan.hellyer.kiwi](https://ryan.hellyer.kiwi) | [GitHub Repo](https://github.com/ryanhellyer/gitmeh)
 
