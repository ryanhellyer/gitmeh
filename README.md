# gitmeh ¯\_(ツ)_/¯

**AI-powered git commits for the terminally lazy.**

**gitmeh** is a high-speed shortcut for your personal garbage repositories. It is designed specifically for those projects where quality does not matter and the only thing you care about is closing the laptop as fast as humanly possible.

> **⚠️ WARNING:** Using this on a professional team project is a great way to get a stern talking-to from a Senior Dev. This tool is reckless, indifferent, and definitely not "enterprise-ready."

![gitmeh in action](images/screenshot.avif)

### Why use this?

Because writing thoughtful commit messages for your 14th unfinished side project is a waste of your precious nap time.

* **Nuclear Staging:** It runs `git add --all` without asking. It stages your unfinished thoughts, your secrets, and that one large `test.mp4` you forgot was there.
* **AI Guesswork:** It begs Google’s Gemini to explain what you did because you have already forgotten.
* **Automatic Pushing:** Shovels your changes directly to the cloud so you can stop looking at the terminal.
* **Built-in Judgement:** Features 40+ randomized status messages that mock your lack of professional standards.

### Quick Start

1. **Get a Gemini API Key** from [Google AI Studio](https://aistudio.google.com/).
2. **Dump it in your shell config** (`~/.bashrc` or `~/.zshrc`):
`export GEMINI_API_KEY='your_key_here'`
3. **Install the thing globally** so you can run it from anywhere without that annoying `.sh` extension:
```bash
mv gitmeh.sh gitmeh
chmod +x gitmeh
sudo mv gitmeh /usr/local/bin/
```

### Requirements

* `git`: duh!
* `jq`: to handle the robot's feelings.
* `curl`: to send the SOS signal to Google.

### Author

**Ryan Hellyer** [ryan.hellyer.kiwi](https://ryan.hellyer.kiwi) | [GitHub Repo](https://github.com/ryanhellyer/gitmeh)
