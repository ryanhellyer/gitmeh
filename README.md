# gitmeh ¯\_(ツ)_/¯

#####################################
#####################################

Git meh API fixes required

** The API should detect errors, and if it finds an error from the API, it should immediately attempt to use another Open Router model. If it does need to serve an error after that, then it should be a consistent error that git meh can always recognise and serve a relevant error message for. **


# I THINK THIS MAY BE CAUSED BY THE BODY BEING TOO LARGE - it happened when I added a bunch of binaries
502 Bad Gateway | raw body: "OpenRouter error: Provider returned error"


❯ git meh

Suggested commit message:
```git

Accept this message? [Y]es / [n]o / [e]dit: 



Go-test on  test [?] via 🐹 
❯ git meh
502 Bad Gateway | raw body: "OpenRouter error: Provider returned error"

Go-test on  test [+] via 🐹 
❯ git meh
Post "https://ai.hellyer.kiwi/gitmeh": context deadline exceeded (Client.Timeout exceeded while awaiting headers)

Go-test on  test [+] via 🐹 

#####################################
#####################################


**AI-powered git commits for the terminally lazy.**

**gitmeh** is a high-speed shortcut for your personal garbage repositories. It is designed specifically for those projects where quality does not matter and the only thing you care about is closing the laptop as fast as humanly possible.

> **⚠️ WARNING:** Using this on a professional team project is a great way to get a stern talking-to from your engineering manager. This tool is reckless, indifferent, and definitely not "enterprise-ready."

![gitmeh in action](images/screenshot.avif)

### Why use this?

Because writing thoughtful commit messages for your 14th unfinished side project is a waste of your precious nap time.

* **Nuclear Staging:** It runs `git add --all` without asking. It stages your unfinished thoughts, your secrets, and that one large `test.mp4` you forgot was there.
* **AI Guesswork:** By default the tool calls an **OpenAI-compatible** chat API (default base `https://openrouter.ai/api/v1`) using a public bearer set in `internal/config.DefaultPublicAPIKey` when non-empty. Set `GITMEH_API_KEY` or `OPENROUTER_API_KEY` to use your own key and optionally `GITMEH_API_BASE` for another compatible host.
* **Automatic Pushing:** Shovels your changes directly to the cloud so you can stop looking at the terminal.
* **Built-in Judgement:** Features 40+ randomized status messages that mock your lack of professional standards.

### Quick Start

1. **Default (no env API key):** The tool POSTs JSON to **`/v1/chat/completions`** on the configured API root (default OpenRouter) using the compiled-in public bearer in `internal/config` when set there. Until that constant is populated, set **`OPENROUTER_API_KEY`** (or **`GITMEH_API_KEY`**) as in step 3. To smoke-test an endpoint, run `./scripts/verify-openai-chat.sh` with **`OPENROUTER_API_KEY`** or **`GITMEH_VERIFY_API_KEY`** set (expect HTTP 200 and a non-empty commit line).

2. **Legacy plain POST (opt-in):** Set `GITMEH_LEGACY_PLAIN=true` to use the old **`text/plain`** flow against `GITMEH_DEFAULT_URL` (default `https://ai.hellyer.kiwi/gitmeh`).

3. **Optional — your own API key:** **Get an OpenRouter API key** from [OpenRouter](https://openrouter.ai/keys) and **dump it in your shell config** (`~/.bashrc` or `~/.zshrc`):

   ```bash
   export OPENROUTER_API_KEY='your_key_here'
   ```

   With that set, **`git meh`** uses OpenRouter (or set `GITMEH_API_BASE` to another OpenAI-compatible root). Optional: `OPENROUTER_MODEL` / `GITMEH_MODEL` (default on OpenRouter: `google/gemma-3-4b-it`). See [openrouter.ai/models](https://openrouter.ai/models).  
   Optional: `GITMEH_PROMPT` to customize the system instructions (the unified diff is always a separate user message).

4. **Install** the `git-meh` binary on your `PATH` (see below). Git discovers it as a subcommand, so you run **`git meh`** from any repository.

### Install

**macOS / Linux** — from the repository root:

```bash
./install.sh
```

That installs into `~/.local/bin` and updates your shell config so that directory is on your `PATH` when needed. Use a new terminal window, or run the `source …` command the script prints, then **`git meh`**.

**Windows:** Put **`git-meh.exe`** on your `PATH`, then run **`git meh`**.

### Requirements

* **Git:** duh!
* **Network:** the binary uses Go’s HTTP client (no shell `curl` required for the default path).

### Unit Tests

From the repository root, with Go installed:

```bash
go vet ./...
go test ./... -count=1
```

### Integration Tests

Integration tests use `git`. They need **`git` installed on your system **

```bash
go test -tags=integration ./... -count=1
```

### Changelog

* **`3.0`:** Rewrite in Go; install the `git-meh` binary and run **`git meh`** (the old shell `gitmeh` command is gone).
* **`3.x`:** Default path uses OpenRouter-compatible chat with optional compiled-in public API key; legacy `text/plain` via `GITMEH_LEGACY_PLAIN=true`.
* `2.1.0`: Default to the free hosted plain-text API so you can avoid another signup; OpenRouter when you set `OPENROUTER_API_KEY`; whine about the 1000 requests/day/IP limit on the free tier
* `2.0.2`: Fixing default model documentation
* `2.0.1`: Set default model to Google Gemma 3 4B as it is free
* `2.0`: Conversion to use OpenRouter API and implementing ability to change model used and prompt
* `1.0`: Initial implementation using Google Gemini

### Author

**Ryan Hellyer** [ryan.hellyer.kiwi](https://ryan.hellyer.kiwi) | [GitHub Repo](https://github.com/ryanhellyer/gitmeh)
 
