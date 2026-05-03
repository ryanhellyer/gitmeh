#!/usr/bin/env bash
# Verify an OpenAI-compatible /v1/chat/completions endpoint (e.g. OpenRouter).
# Usage: export OPENROUTER_API_KEY=...  OR  GITMEH_VERIFY_API_KEY=...
#         ./scripts/verify-openai-chat.sh
# Optional: GITMEH_VERIFY_BASE (default https://openrouter.ai/api/v1),
#           GITMEH_VERIFY_MODEL (default google/gemma-3-4b-it).

set -euo pipefail

KEY="${GITMEH_VERIFY_API_KEY:-${OPENROUTER_API_KEY:-}}"
if [[ -z "${KEY}" ]]; then
	echo "verify-openai-chat: set GITMEH_VERIFY_API_KEY or OPENROUTER_API_KEY" >&2
	exit 1
fi

BASE="${GITMEH_VERIFY_BASE:-https://openrouter.ai/api/v1}"
BASE="${BASE%/}"
MODEL="${GITMEH_VERIFY_MODEL:-google/gemma-3-4b-it}"
URL="$BASE/chat/completions"
export GITMEH_VERIFY_MODEL_FOR_BODY="$MODEL"

if ! command -v python3 >/dev/null 2>&1; then
	echo "python3 is required for this script" >&2
	exit 1
fi

BODY="$(python3 <<'PY'
import json
import os

sys_prompt = "You write git commit messages. Reply with one line only, no preamble."
diff = """Unified diff:
--- a/foo
+++ b/foo
@@ -1 +1 @@
-x
+y
"""
model = os.environ["GITMEH_VERIFY_MODEL_FOR_BODY"]
print(json.dumps({
    "model": model,
    "messages": [
        {"role": "system", "content": sys_prompt},
        {"role": "user", "content": diff},
    ],
    "temperature": 0.3,
    "max_tokens": 64,
}))
PY
)"

tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT

code="$(curl -sS -o "$tmp" -w "%{http_code}" \
	-X POST "$URL" \
	-H "Content-Type: application/json" \
	-H "Accept: application/json" \
	-H "Authorization: Bearer ${KEY}" \
	-d "$BODY")"

if [[ "$code" != "200" ]]; then
	echo "verify-openai-chat: expected HTTP 200 from $URL, got $code" >&2
	cat "$tmp" >&2 || true
	exit 1
fi

python3 <<PY
import json, sys
path = "$tmp"
with open(path) as f:
    data = json.load(f)
choices = data.get("choices") or []
if not choices:
    print("verify-openai-chat: no choices in response", file=sys.stderr)
    sys.exit(1)
content = (choices[0].get("message") or {}).get("content")
if not (isinstance(content, str) and content.strip()):
    print("verify-openai-chat: empty choices[0].message.content", file=sys.stderr)
    sys.exit(1)
print(content.strip())
PY

echo "verify-openai-chat: OK"
