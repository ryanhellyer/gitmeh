# gitmeh ¯\_(ツ)_/¯

**AI-powered git commits for the terminally lazy.**

**gitmeh** is a high-speed shortcut for your personal garbage repositories. It is designed specifically for those projects where quality does not matter and the only thing you care about is closing the laptop as fast as humanly possible.

> **⚠️ WARNING:** Using this on a professional team project is a great way to get a stern talking-to from your engineering manager. This tool is reckless, indifferent, and definitely not "enterprise-ready."

![gitmeh in action](images/screenshot.avif)

### Why use this?

Because writing thoughtful commit messages for your 14th unfinished side project is a waste of your precious nap time.

* **Nuclear Staging:** It runs `git add --all` without asking. It stages your unfinished thoughts, your secrets, and that one large `test.mp4` you forgot was there.
* **AI Guesswork:** It begs OpenRouter (any supported model) to explain what you did because you have already forgotten.
* **Automatic Pushing:** Shovels your changes directly to the cloud so you can stop looking at the terminal.
* **Built-in Judgement:** Features 40+ randomized status messages that mock your lack of professional standards.

### Quick Start

1. **Get an OpenRouter API key** from [OpenRouter](https://openrouter.ai/keys).
2. **Dump it in your shell config** (`~/.bashrc` or `~/.zshrc`):
   ```bash
   export OPENROUTER_API_KEY='your_key_here'
   ```
   Optional: set `OPENROUTER_MODEL` (default: `google/gemma-3-4b-it`). See [openrouter.ai/models](https://openrouter.ai/models).  
   Optional: set `GITMEH_PROMPT` to customize the instruction sent to the AI (the diff is always appended).
3. **Install the thing globally** so you can run it from anywhere without that annoying `.sh` extension:

macOS / Linux:
```bash
chmod +x gitmeh.sh && sudo mv gitmeh.sh /usr/local/bin/gitmeh
```

Windows (Git Bash - _totally untested as I don't use Windows_):
```bash
mkdir -p ~/bin
cp gitmeh.sh ~/bin/gitmeh
# Ensure ~/bin is in your PATH
```

### Requirements

* `git`: duh!
* `jq`: to handle the robot's feelings.
* `curl`: to send the SOS signal to OpenRouter.

### Changelog

* `2.0.2`: Fixing default model documentation
* `2.0.1`: Set default model to Google Gemma 3 4B as it is free
* `2.0`: Conversion to use OpenRouter API and implementing ability to change model used and prompt
* `1.0`: Initial implementation using Google Gemini

### Author

**Ryan Hellyer** [ryan.hellyer.kiwi](https://ryan.hellyer.kiwi) | [GitHub Repo](https://github.com/ryanhellyer/gitmeh)
