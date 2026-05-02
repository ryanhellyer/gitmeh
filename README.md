# gitmeh ¯\_(ツ)_/¯

#####################################
#####################################

Git meh API fixes required

** The API should detect errors, and if it finds an error from the API, it should immediately attempt to use another Open Router model. If it does need to serve an error after that, then it should be a consistent error that git meh can always recognise and serve a relevant error message for. **


# I THINK THIS MAY BE CAUSED BY THE BODY BEING TOO LARGE - it happened when I added a bunch of binaries
502 Bad Gateway | raw body: "OpenRouter error: Provider returned error"


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
* **AI Guesswork:** By default it flings your staged diff at a free hosted API so you do not have to pretend you will ever sign up for anything. It still explains what you did because you have already forgotten. If you are picky, wave an OpenRouter key around and make it beg a model of your choosing instead.
* **Automatic Pushing:** Shovels your changes directly to the cloud so you can stop looking at the terminal.
* **Built-in Judgement:** Features 40+ randomized status messages that mock your lack of professional standards.

### Quick Start

1. **Default (no API key):** Do nothing special. The tool POSTs your **staged diff** to `https://ai.hellyer.kiwi/gitmeh` as `text/plain` and uses the response as the commit message. The free tier is **limited to 1000 requests per day per IP address**, so if you and your twelve roommates all commit-spam at once, you will hit the wall together. Pace yourselves.

> THE OPENROUTER CODE STILL NEEDS TO BE INCLUDED IN THE GO REWRITE.

2. **Optional — OpenRouter:** If you insist on owning the relationship, **get an OpenRouter API key** from [OpenRouter](https://openrouter.ai/keys) and **dump it in your shell config** (`~/.bashrc` or `~/.zshrc`):

   ```bash
   export OPENROUTER_API_KEY='your_key_here'
   ```

   With that set, **`git meh`** bothers OpenRouter instead of the default URL. Optional: `OPENROUTER_MODEL` (default: `google/gemma-3-4b-it`). See [openrouter.ai/models](https://openrouter.ai/models).  
   Optional: `GITMEH_PROMPT` to customize the instruction sent to the AI (the diff is always appended; OpenRouter only — the free endpoint does not care about your feelings).  
   Optional: `GITMEH_DEFAULT_URL` if you want a different keyless endpoint (full URL; default: `https://ai.hellyer.kiwi/gitmeh`).

3. **Install** the `git-meh` binary on your `PATH` (see below). Git discovers it as a subcommand, so you run **`git meh`** from any repository.

### Install

**macOS / Linux** — from the repository root:

```bash
./install.sh
```

That installs into `~/.local/bin` and updates your shell config so that directory is on your `PATH` when needed. Use a new terminal window, or run the `source …` command the script prints, then **`git meh`**.

**Windows:** Put **`git-meh.exe`** on your `PATH`, then run **`git meh`**.

### Requirements

* **Git:** duh!
* **`curl`:** to send the SOS signal (default API or OpenRouter).
* **`jq`:** to handle the robot's feelings — **only if** you are on the OpenRouter path. The keyless mode does not need it, because apparently plain text is easier than JSON.

### Changelog

* **`3.0`:** Rewrite in Go; install the `git-meh` binary and run **`git meh`** (the old shell `gitmeh` command is gone).
* `2.1.0`: Default to the free hosted plain-text API so you can avoid another signup; OpenRouter when you set `OPENROUTER_API_KEY`; whine about the 1000 requests/day/IP limit on the free tier
* `2.0.2`: Fixing default model documentation
* `2.0.1`: Set default model to Google Gemma 3 4B as it is free
* `2.0`: Conversion to use OpenRouter API and implementing ability to change model used and prompt
* `1.0`: Initial implementation using Google Gemini

### Author

**Ryan Hellyer** [ryan.hellyer.kiwi](https://ryan.hellyer.kiwi) | [GitHub Repo](https://github.com/ryanhellyer/gitmeh)
 
