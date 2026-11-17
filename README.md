# gitmeh ¯\_(ツ)_/¯

# todo

1. Make sure response is formatted correctly, perhaps can check for all the possible Conventional Commits starts to confirm.
2. Add attribute --auto / -a (automatically accepts commit message
3. Add attribute -- ??
4. Handle huge diffs without erroring

-----

**AI-powered git commits for people with better things to do.**

Stages all changes, generates a commit message via AI, commits, and pushes — all in one command. Review and edit before it pushes.

## Quick Start

```bash
# 1. Install
curl -fsSL https://raw.githubusercontent.com/ryanhellyer/gitmeh/master/install.sh | bash

# 2. Run
git meh
```

Or build from source: `make build && cp git-meh ~/.local/bin/` (requires Go).

Git discovers the binary as a subcommand — works in any repository. No API key required — gitmeh ships with a built-in default that works out of the box.

> **⚠️** Review what you're staging before AI pushes it, especially on shared repos.

### What it does

Stages everything (`git add --all`), generates a commit message via AI, commits it, and pushes. For those who believe commit messages are important — just not important enough to write by hand.

Because commit messages matter — just not enough to spend brainpower on them.

* **Automated Staging:** Runs `git add --all` so you don't have to think about what changed.
* **AI-Generated Messages:** Produces a commit message via the gitmeh API (`POST /v1/git_meh/completions`) or any OpenAI-compatible chat API, with retry logic, fallback models, and exponential backoff.
* **Interactive Review:** Review the message before committing, edit it inline, or abort — all at a single prompt.
* **Automatic Pushing:** Commits and pushes in one step.

### Default API service

If you don't set `GITMEH_API_KEY`, gitmeh uses a **free hosted API** at `https://ai.hellyer.kiwi/v1`, run by the author (Ryan Hellyer). The client POSTs the staged diff to `https://ai.hellyer.kiwi/v1/git_meh/completions`, and the server applies the standard commit-message instructions (Conventional Commits format, length, tone) before calling the model. Keeping the instructions server-side means every client — the CLI, a script, any bot — gets the same output contract. The backend is a Laravel-based API and currently uses **Ling 3.0 Flash**, which produces creditably good commit messages for a tool that started as a throwaway Bash script. If usage costs climb too high I may need to switch to a smaller model, but for now the quality/price ratio is excellent.

The endpoint contract and the exact instructions are documented in [`docs/git_meh_endpoint.md`](docs/git_meh_endpoint.md).

**Your staged diff (code) is sent to this server** and forwarded to whichever model the backend selects. If you're uncomfortable sharing your code with me (Ryan) or with the third-party model provider, **do not use the default service**. Set `GITMEH_API_BASE`, `GITMEH_API_KEY`, and `GITMEH_MODEL` to use your own AI provider instead.

I have zero interest in your code and no intention of reading it, but it does pass through my server and the model provider's servers.

### Using your own API key (optional)

If you'd prefer to use your own AI provider instead of the default service, set at minimum:

```bash
export GITMEH_API_BASE='https://opencode.ai/zen/v1'   # or any OpenAI-compatible endpoint
export GITMEH_API_KEY='your_api_key'
```

Works with OpenCode Zen, OpenAI, OpenRouter, and any OpenAI-compatible API. With a user key, gitmeh calls `{GITMEH_API_BASE}/chat/completions` and sends the Conventional Commits prompt as the system message; without a key it calls `{GITMEH_API_BASE}/git_meh/completions` and the server supplies that prompt. A custom `GITMEH_API_BASE` with no key is therefore assumed to be a gitmeh-compatible API (for example the staging server), not a bare OpenAI-compatible host.

All available config options when bringing your own key:

| Env var | Description | Default |
| --- | --- | --- |
| `GITMEH_API_BASE` | API base URL | `https://ai.hellyer.kiwi/v1` (built-in) |
| `GITMEH_API_KEY` | API key | built-in public key |
| `GITMEH_MODEL` | Model name | `gitmeh-hosted` or `google/gemma-3-4b-it` |
| `GITMEH_PROMPT` | Instructions override (sent as `instructions` to the gitmeh endpoint, or as the system prompt to a chat endpoint) | server-side standard prompt, or Conventional Commits prompt with your own key |
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

### 3.1 — 2026-11-17
Switched default API to use https://ai.hellyer.kiwi/v1/git_meh/completions.

### 3.0.3 — 2026-07-17
Allow insecure TLS skip for any host in dev mode. This allows for simple testing of the API without requiring https to be configured on the test API.

### 3.0.2 — 2026-06-07

Bumped all dependencies: `golang.org/x/term`, `actions/checkout`, `actions/setup-go`, `golangci/golangci-lint-action`, `codecov/codecov-action`, `softprops/action-gh-release`.

### 3.0.1 — 2026-06-07

Fixed release automation: GitHub Actions now builds and attaches cross-platform binaries (Linux/macOS, amd64/arm64) to releases. `install.sh` downloads these assets — no Go toolchain required to install.

### 3.0 — 2026-06-07

Complete rewrite in Go. Ships with a built-in API key and hosted endpoint — works out of the box with zero configuration. Previous versions required your own OpenRouter API key.

### 2.0.2 — 2026-03-17

Fixed default model selection no longer defaults to a paid model. (Also tagged as 2.0.1 — same release, corrected tag.)

### 2.0 — 2026-03-11

Updated API to an OpenAI-compatible backend.

### 1.0 — 2026-03-11

Initial implementation using Google Gemini.

## Author

**Ryan Hellyer** [ryan.hellyer.kiwi](https://ryan.hellyer.kiwi) | [GitHub Repo](https://github.com/ryanhellyer/gitmeh)
 