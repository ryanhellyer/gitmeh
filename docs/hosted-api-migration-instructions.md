# Hosted gitmeh API: OpenAI-compatible chat (server-side work)

Keep the existing /gitmeh path with it's existing (now legacy) functionality. But add a new endpoint to handle this API in the same format as a typical Open AI compatible API endpoint would work. Currently, only the one existing API format should be approved, which would be the one in the legacy functionality, just as a typical raw Open AI API endpoint rather than it being filtered and given pre-instructions like it is now.



| Item | Value |
|------|--------|
| Public bearer token (weak client id, not a billing secret) | `gitmeh-public-client` |
| Default chat API base URL (no trailing slash) | `https://ai.hellyer.kiwi/v1` |
| Full completions URL | `https://ai.hellyer.kiwi/v1/chat/completions` |
| Legacy plain endpoint (keep during transition) | `POST https://ai.hellyer.kiwi/gitmeh` with `Content-Type: text/plain; charset=UTF-8`, body = raw unified diff, response = plain text commit message |
| Keep legacy plain path | **Yes** until old binaries are gone; same per-IP limits as today |
| Max JSON request body | **2097152** bytes (2 MiB) before parsing; reject larger with `413` or `400` and a short JSON error |

## Reference client behavior (must match)

- **Method/path:** `POST {baseURL}/chat/completions` where `baseURL` has **no trailing slash** (client uses `base + "/chat/completions"`).
- **Headers (request):**
  - `Content-Type: application/json`
  - `Accept: application/json`
  - `Authorization: Bearer <token>` — the published CLI uses `GITMEH_API_KEY` or `OPENROUTER_API_KEY` when set; otherwise a bearer string injected at **link time** (`-X gitmeh/internal/config.BuiltinAPIKey=...` in release builds). Treat the token as **optional** or **required** on the server, but document which; mismatches should return **401** with JSON error if you require it.
- **JSON request body (minimum fields to support):**
  - `model` (string) — client sends the configured model id (default on OpenRouter: `google/gemma-3-4b-it`); you may **ignore** and always run your local model, or **map** ids; return **400** if you require a specific model and it is missing.
  - `messages` (array) — client sends **two** messages: `role: "system"` (instructions) and `role: "user"` with content `Unified diff:\n` + unified diff text.
  - `temperature` (number, e.g. 0.3) — optional to honor; safe to clamp.
  - `max_tokens` (number, e.g. 512) — optional to honor; cap server-side for cost control.
- **JSON response body (success):** OpenAI shape, at minimum:
  - `choices` non-empty array
  - `choices[0].message.content` string = assistant commit message (plain text; client trims whitespace).
- **Errors (non-2xx):** Prefer `{"error":{"message":"...","type":"...","code":"..."}}` so clients can surface `error.message`.

## Auth policy (implemented on server)

Use **optional Bearer** for `https://ai.hellyer.kiwi/v1/chat/completions`:

- If `Authorization: Bearer gitmeh-public-client` matches, treat as official gitmeh client (same rate limits as legacy).
- If header missing or wrong: either same limits (public) or **401** — pick one and document; client always sends the bearer for the hosted default.

## Server implementation checklist

1. **Routing:** `POST` on `/v1/chat/completions` (under your TLS host); **405** for wrong methods.
2. **Parse JSON** with **2 MiB** max body; reject oversized bodies before model call.
3. **Extract diff** from `messages`: prefer the last `user` message; strip an optional `Unified diff:\n` prefix for robustness.
4. **Build prompt** for your local model: system text from `role == "system"` messages; user content = diff.
5. **Reuse** the same inference path as the legacy `text/plain` endpoint so behavior stays consistent.
6. **Response:** `Content-Type: application/json`; **200** with `choices[0].message.content` set to the commit message only (no markdown fences).
7. **Rate limiting:** Same per-IP limits as legacy `/gitmeh`.
8. **Timeouts:** Compatible with ~20s client HTTP timeout.
9. **Logging:** Status, latency; avoid logging full diffs if privacy matters.
10. **Tests:** Happy path JSON; missing `messages`; empty `choices`; malformed JSON; oversize body; rate limit if testable.

## Deployment notes

- Reverse proxy: allow maximum of 0.5 MB upload for this route.
- Preserve real client IP for rate limiting (`X-Forwarded-For` trust).

## After the API is live

The reference CLI in this repo now defaults to a generic OpenAI-compatible host; to verify any `/v1/chat/completions` deployment (including staging for this service), run `./scripts/verify-openai-chat.sh` with `OPENROUTER_API_KEY` or `GITMEH_VERIFY_API_KEY` set, and optional `GITMEH_VERIFY_BASE` / `GITMEH_VERIFY_MODEL`. Expect HTTP **200** and non-empty `choices[0].message.content`.
