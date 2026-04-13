#!/usr/bin/env bash

# gitmeh: AI-powered git commits for the terminally lazy.
# Author: Ryan Hellyer <ryan@hellyer.kiwi>
# Website: https://ryan.hellyer.kiwi
# GitHub: https://github.com/ryanhellyer/gitmeh

# Configuration
# Default: POST plain-text diff to the free hosted API (no key). Override with OpenRouter by setting OPENROUTER_API_KEY.
GITMEH_DEFAULT_URL="${GITMEH_DEFAULT_URL:-https://ai.hellyer.kiwi/gitmeh}"
API_KEY="${OPENROUTER_API_KEY:-$GEMINI_API_KEY}"
MODEL="${OPENROUTER_MODEL:-google/gemma-3-4b-it}"
BRANCH=$(git rev-parse --abbrev-ref HEAD)
MAX_TOTAL_CHARS=10000
CHARS_PER_FILE=800

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' 

# --- The "Lazy" Phrase Arrays ---

INTRO_PHRASES=(
    "gitmeh: For when your career is a series of shortcuts."
    "gitmeh: Lowering the bar for commit history since today."
    "gitmeh: Because typing is the enemy."
    "gitmeh: The 'I am just here for the paycheck' utility."
    "gitmeh: For the developer who has truly given up."
    "gitmeh: Automating your lack of interest."
    "gitmeh: Because 'fixed stuff' isn't a professional commit message."
    "gitmeh: Helping you pretend you worked hard today."
    "gitmeh: Your personal ghostwriter for mediocrity."
    "gitmeh: The 'close the laptop and walk away' button."
)

STAGING_PHRASES=(
    "Staging everything because you're too lazy to pick..."
    "Tossing everything into the stage like a laundry pile..."
    "Adding everything because nuance is for people with energy..."
    "Staging your messy room of code... don't look too closely."
    "Grabbing everything. Hope there's no secrets in there. (There probably are.)"
    "Nuclear staging initiated. RIP clean history."
    "Shoveling your changes into the commit bucket..."
    "Staging everything. It is not like you were going to review it anyway."
    "Blindly adding files because life is too short for git add -p."
    "Bulk staging. Let God (or the AI) sort them out."
)

THINKING_PHRASES=(
    "Consulting the robot because thinking is hard..."
    "Asking the AI to lie about how much work you did..."
    "Delegating your cognitive load to a server in Oregon..."
    "Letting the algorithm pretend you are a professional..."
    "Begging the AI to explain your own code back to you..."
    "Outsourcing your last two brain cells to the cloud..."
    "Waiting for the robot to find a nice way to say 'you broke it'..."
    "Requesting a miracle from the API..."
    "Pinging the mothership for a crumb of inspiration..."
    "Asking the AI to cover for you. Again."
)

SUCCESS_PHRASES=(
    "It's pushed. Go outside."
    "The deed is done. Go be useless elsewhere."
    "Done. Don't check the logs. Just walk away."
    "Success. Your secret is safe with the AI."
    "It's in the cloud now. Not your problem anymore."
    "Mission accomplished. Nap time."
    "Pushed. Let the CI/CD pipeline deal with your mess."
    "And... stay out. See you tomorrow (maybe)."
    "Finished. That's enough 'work' for one day."
    "The code is gone. Fly, little bird, fly."
)

# --- Helper to pick a random phrase ---
get_random() {
    local arr=("$@")
    echo "${arr[$(( RANDOM % ${#arr[@]} ))]}"
}

# --- Script Logic ---

# Help instructions
if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
    echo -e "${CYAN}$(get_random "${INTRO_PHRASES[@]}")${NC}"
    echo "Usage: gitmeh"
    echo ""
    echo "By default this uses the free hosted API (no key). Optional: use your own OpenRouter account:"
    echo "  export OPENROUTER_API_KEY='your_key_here'"
    echo ""
    echo "Optional (OpenRouter): OPENROUTER_MODEL (default: google/gemma-3-4b-it). See https://openrouter.ai/models"
    echo "Optional (OpenRouter): GITMEH_PROMPT to customize the instruction (the diff is always appended)."
    echo "Optional: GITMEH_DEFAULT_URL for a different keyless endpoint URL (default: https://ai.hellyer.kiwi/gitmeh)"
    echo ""
    echo "Author: Ryan Hellyer (https://ryan.hellyer.kiwi)"
    exit 0
fi

# jq is only required when using OpenRouter (JSON responses)
if [ -n "$API_KEY" ] && ! command -v jq &> /dev/null; then
    echo "Error: 'jq' is missing. Install it, or unset OPENROUTER_API_KEY to use the keyless default API."
    exit 1
fi

# Add changes
git add --all
echo -e "${CYAN}$(get_random "${STAGING_PHRASES[@]}")${NC}"
git status --short

FILES=$(git diff --cached --name-only)

if [ -z "$FILES" ]; then
    echo "No changes. Go back to sleep."
    exit 0
fi

echo -e "\n$(get_random "${THINKING_PHRASES[@]}")"

if [ -n "$API_KEY" ]; then
    # OpenRouter: chunked diff + JSON API
    SMART_DIFF=""
    for FILE in $FILES; do
        if [ ${#SMART_DIFF} -gt $MAX_TOTAL_CHARS ]; then
            SMART_DIFF+=$'\n' "... [Truncated because I'm bored] ..."
            break
        fi
        FILE_DIFF=$(git diff --cached -- "$FILE" | head -c $CHARS_PER_FILE)
        SMART_DIFF+=$'\n'"--- File: $FILE ---"$'\n'"$FILE_DIFF"$'\n'
    done

    GITMEH_PROMPT_DEFAULT="Write a short, professional git commit message for these changes. Use imperative mood. Only return the message text:"
    PROMPT="${GITMEH_PROMPT:-$GITMEH_PROMPT_DEFAULT} $SMART_DIFF"
    JSON_PAYLOAD=$(jq -n --arg msg "$PROMPT" --arg model "$MODEL" '{model: $model, messages: [{role: "user", content: $msg}]}')

    RESPONSE=$(curl -sS -X POST "https://openrouter.ai/api/v1/chat/completions" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_KEY" \
        -d "$JSON_PAYLOAD")

    API_ERR=$(echo "$RESPONSE" | jq -r '.error.message // .error // empty' 2>/dev/null)
    if [ -n "$API_ERR" ]; then
        echo -e "${YELLOW}OpenRouter error: ${API_ERR}${NC}"
        exit 1
    fi

    COMMIT_MSG=$(echo "$RESPONSE" | jq -r '.choices[0].message.content // empty' 2>/dev/null | sed 's/^[[:space:]]*//;s/[[:space:]]*$//' | head -1)
else
    # Default: plain-text cached diff to hosted API (prompt is server-side)
    DIFF_BODY=$(git diff --cached)
    if [ ${#DIFF_BODY} -gt $MAX_TOTAL_CHARS ]; then
        DIFF_BODY="${DIFF_BODY:0:$MAX_TOTAL_CHARS}"$'\n'"... [Truncated because I'm bored] ..."
    fi

    if ! RESPONSE=$(curl -sS -f -X POST "$GITMEH_DEFAULT_URL" \
        -H "Content-Type: text/plain" \
        --data-binary "$DIFF_BODY"); then
        echo -e "${YELLOW}Request to default API failed (network or HTTP error). Set OPENROUTER_API_KEY to use OpenRouter instead.${NC}"
        exit 1
    fi

    COMMIT_MSG=$(printf '%s' "$RESPONSE" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//' | head -1)
fi

if [ -z "$COMMIT_MSG" ] || [ "$COMMIT_MSG" == "null" ]; then
    echo -e "${YELLOW}The AI failed. Probably went on a coffee break.${NC}"
    echo "Response snippet: $(echo "$RESPONSE" | head -c 500)"
    exit 1
fi

# Confirmation
echo -e "------------------------------------------------"
echo -e "Proposed: ${GREEN}${COMMIT_MSG}${NC}"
echo -e "------------------------------------------------"
read -p "Commit and push? [Y/n/e]: " USER_INPUT
USER_INPUT=${USER_INPUT:-y}

case "$USER_INPUT" in
    [yY][eE][sS]|[yY]) 
        git commit -m "$COMMIT_MSG"
        if git push origin "$BRANCH"; then
            echo -e "${CYAN}$(get_random "${SUCCESS_PHRASES[@]}")${NC}"
        else
            echo -e "${YELLOW}Push failed. You actually have to do some work now (git pull).${NC}"
        fi
        ;;
    [eE][dD][iI][tT]|[eE])
        read -p "Fine, fix it yourself: " MANUAL_MSG
        git commit -m "$MANUAL_MSG"
        git push origin "$BRANCH"
        echo -e "${CYAN}Manually fixed and pushed. Look at you go.${NC}"
        ;;
    *)
        echo "Aborted. Coward."
        exit 1
        ;;
esac