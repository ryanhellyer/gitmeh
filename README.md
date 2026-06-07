# gitmeh ¯\_(ツ)_/¯

**AI-powered git commits for the terminally lazy.**

Stages everything (`git add --all`), AI-guesses a commit message, then shovels it to the cloud. Designed for people who can't be bothered writing their own commit messages.

This started life as a Bash script written on a Sunday afternoon. Then people actually started using it, so I rewrote it in Go — partly as a learning exercise, partly because the Bash version had too many sharp edges. The result is faster, more reliable, and genuinely useful enough that I've softened my stance on whether it belongs in professional workflows.

> **⚠️** Review what you're staging before you let AI push it, especially on shared repos. That said, the "stage, AI-draft, review, commit" workflow is genuinely faster than writing messages by hand — I use it daily now.

### Why use this?

Because writing commit messages takes effort and you've got better things to do.

* **Automated Staging:** Runs `git add --all` so you don't have to think about what changed.
* **AI Guesswork:** Generates a commit message via an OpenAI-compatible chat API, with retry logic, fallback models, and exponential backoff.
* **Interactive Review:** Review the message before committing, edit it inline with cursor keys, or abort — all at a single prompt.
* **Automatic Pushing:** Commits and pushes in one step.

### Default API service

If you don't set `GITMEH_API_KEY`, gitmeh uses a **free hosted API** at `https://ai.hellyer.kiwi/v1`, run by the author (Ryan Hellyer). The backend is a Laravel-based API and currently runs **Deepseek V4 Flash**, which produces surprisingly good commit messages for a tool that started life as a Bash joke. If usage costs climb too high I may need to switch to a smaller model, but for now the quality/price ratio is excellent.

**Your staged diff (code) is sent to this server** and then forwarded to whichever model the backend picks. If you are not comfortable sharing your code with me (Ryan) or with the random third-party model I route it through, **do not use the default service**. Set `GITMEH_API_BASE`, `GITMEH_API_KEY`, `GITMEH_MODEL` etc. to use your own AI provider instead.

I have zero interest in your code and no intention of looking at it, but it will be processed through my server and the model provider's servers.

## Quick Start

```bash
# 1. Install
make build && cp git-meh ~/.local/bin/           # from the repo root (requires Go)
# Or: ./install.sh                              # uses a prebuilt binary

# 2. Run
git meh
```

Git discovers the binary as a subcommand — works in any repository. No API key required — gitmeh ships with a built-in default that works out of the box.

### Using your own API key (optional)

If you'd prefer to use your own AI provider instead of the default service, set at minimum:

```bash
export GITMEH_API_BASE='https://opencode.ai/zen/v1'   # or any OpenAI-compatible endpoint
export GITMEH_API_KEY='your_api_key'
```

Works with OpenCode Zen, OpenAI, OpenRouter, and any OpenAI-compatible API.

All available config options when bringing your own key:

| Env var | Description | Default |
|---|---|---|---|
| `GITMEH_API_BASE` | API base URL | `https://ai.hellyer.kiwi/v1` (built-in) |
| `GITMEH_API_KEY` | API key | built-in public key |
| `GITMEH_MODEL` | Model name | `gitmeh-hosted` or `google/gemma-3-4b-it` |
| `GITMEH_PROMPT` | System prompt for the model | Conventional Commits prompt |
| `GITMEH_FALLBACK_MODELS` | Comma-separated models to try if the primary fails | — |
| `GITMEH_MAX_DIFF_BYTES` | Per-file diff truncation limit (0 = no limit) | `10000` (10 KB) |

**Auth priority**: `GITMEH_API_KEY` > built-in public key.

**Fallback models**: If the primary model fails (timeout, 5xx response codes, context-length exceeded), gitmeh retries up to 3 times with exponential backoff, then tries each fallback model in order. A 401 or other client error skips retries immediately.

**Diff truncation**: When the staged diff exceeds `GITMEH_MAX_DIFF_BYTES`, gitmeh keeps all file headers and proportionally trims hunk content per file. Truncated sections are marked with `# hunk truncated`.

## Developer Guide

### Prerequisites

- Go (see `go.mod` for version)
- `golangci-lint` and `govulncheck` for linting (install via `go install`)

### Commands

```bash
make dev         # build native binary (developer mode, self-signed TLS)
make build       # build native binary (production)
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

### Dev / prod builds

When built with `make dev`, the binary targets `ai.hellyer.test` and accepts self-signed TLS certificates (developer mode). This is controlled by a linker flag (`-ldflags="-X gitmeh/internal/config.isDev=true"`) so the dev hostname is never compiled into release binaries.

`make build` and `make cross` produce release binaries that target `ai.hellyer.kiwi` with full TLS verification.

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
 
