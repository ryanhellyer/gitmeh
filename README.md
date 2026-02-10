# gitmeh ¯\_(ツ)_/¯

**AI-powered git commits for the terminally lazy.**

**gitmeh** is a high-speed shortcut for your personal garbage repos. It’s designed specifically for those projects where quality doesn't matter, and the only thing you care about is closing the laptop as fast as humanly possible.

> **⚠️ WARNING:** Using this on a professional team project is a great way to get a stern talking-to from a Senior>


### Quick Start

1. **Get a Gemini API Key** from Google AI Studio.
2. **Add to shell config** (`~/.bashrc` or `~/.zshrc`):
`export GEMINI_API_KEY='your_key_here'`
3. **Make it executable:** `chmod +x gitmeh`
4. **Run it:** `./gitmeh`

### What it does (to save you effort)

* **Stages everything:** Runs `git add --all` automatically.
* **Consults the robot:** Generates a message based on your mess.
* **Pushes for you:** Commits and pushes in one step.
* **Randomized insults:** Includes 40+ phrases to mock your work ethic.

### Requirements

* `git`
* `jq`
* `curl`

### Author

**Ryan Hellyer**
[ryan.hellyer.kiwi](https://ryan.hellyer.kiwi) | [GitHub Repo](https://github.com/ryanhellyer/gitmeh)
