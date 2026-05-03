#!/usr/bin/env bash
# Verify the hosted gitmeh OpenAI-compatible chat endpoint.
# Usage: ./scripts/verify-hosted-api.sh
# Optional: GITMEH_VERIFY_BASE (default https://ai.hellyer.kiwi/v1),
#           GITMEH_VERIFY_TOKEN (default gitmeh-public-client).

set -euo pipefail

BASE="${GITMEH_VERIFY_BASE:-https://ai.hellyer.kiwi/v1}"
BASE="${BASE%/}"
TOKEN="${GITMEH_VERIFY_TOKEN:-gitmeh-public-client}"
URL="$BASE/chat/completions"

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required for this script" >&2
  exit 1
fi

BODY="$(python3 <<'PY'
import json
sys_prompt = "You write git commit messages. Reply with one line only, no preamble."
diff = """Unified diff:
--- a/foo
+++ b/foo
@@ -1 +1 @@
-x
+y
"""
print(json.dumps({
    "model": "gitmeh-hosted",
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
  -H "Authorization: Bearer ${TOKEN}" \
  -d "$BODY")"

if [[ "$code" != "200" ]]; then
  echo "verify-hosted-api: expected HTTP 200 from $URL, got $code" >&2
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
    print("verify-hosted-api: no choices in response", file=sys.stderr)
    sys.exit(1)
content = (choices[0].get("message") or {}).get("content")
if not (isinstance(content, str) and content.strip()):
    print("verify-hosted-api: empty choices[0].message.content", file=sys.stderr)
    sys.exit(1)
print(content.strip())
PY

echo "verify-hosted-api: OK"
