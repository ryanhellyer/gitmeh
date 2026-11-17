# `POST /v1/git_meh/completions`

gitmeh's own commit-message endpoint. It exists so that the client, and any
other bot or script that knows the URL, gets the same output contract: the
standard Conventional Commits instructions, sampling parameters and response
shape all live here rather than in each caller.

## Contract

| Item | Value |
| --- | --- |
| Method and path | `POST /v1/git_meh/completions` |
| Auth | `Authorization: Bearer <token>` (the built-in public key is `gitmeh-public-client`) |
| Content type | `application/json` |
| Timeout expected by the client | 60 s |

### Request

```json
{
  "model": "gitmeh-hosted",
  "diff": "Unified diff:\ndiff --git a/main.go b/main.go\n...",
  "instructions": "optional override"
}
```

- `diff` — required. The client already prefixes it with `Unified diff:\n`, so
  the endpoint can pass it straight through as the user message.
- `model` — optional. Clients send `gitmeh-hosted` by default and use this
  field for fallback models; map unknown values to your configured upstream
  model instead of rejecting them.
- `instructions` — optional. Sent only when the user sets `GITMEH_PROMPT`.
  When present it **replaces** the standard instructions below.

### Response

Return OpenAI-compatible JSON so existing client parsing keeps working:

```json
{"choices":[{"message":{"role":"assistant","content":"feat(api): add git_meh endpoint"}}]}
```

Errors must be non-2xx with a message the client can surface:

```json
{"error":{"message":"upstream error"}}
```

The client retries `429`, `5xx`, timeouts and connection failures up to three
times per model, then tries the next fallback model. A `4xx` other than `429`
is treated as fatal, so return `500` for upstream trouble and reserve `4xx` for
bad requests.

## Server behaviour

1. Validate `diff` is present and non-empty; reject otherwise with `422`.
2. Build the upstream chat request:
   - system message: `instructions` if supplied, otherwise the standard
     instructions below;
   - user message: `diff` verbatim;
   - `temperature: 0.3`, `max_tokens: 4096`.
3. Trim the returned content, and reject/retry a reply that does not match
   `^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-z0-9._/-]+\))?!?: .+`
   (strip surrounding quotes or code fences before validating).
4. Wrap the result in the `choices` response above.

## Standard instructions

Use this text verbatim as the system message. It is the single source of truth
for format, length, scope and tone; do not let callers weaken it.

```text
You are gitmeh, the commit-message generator behind the `git meh` CLI. You
receive a unified diff of staged changes and reply with exactly one Git commit
message.

WHAT THE MESSAGE IS FOR
The message is committed verbatim with `git commit -m` and pushed. A human
reads it in `git log` and in code review months later. It must say what
changed and why, in the imperative mood, and match the repository's existing
style.

FORMAT - Conventional Commits 1.0.0

<type>[optional scope][!]: <description>

- <type> is one of: feat, fix, docs, style, refactor, perf, test, build, ci,
  chore, revert.
- <scope> is optional, lowercase, and names the affected area (for example
  aiapi, config, git, cli). Use it only when it adds clarity.
- <description> starts lowercase, has no trailing period, and uses the
  imperative mood: "add", "fix", "remove" - never "added", "adds" or
  "adding".
- Mark a breaking change with `!` before the colon, for example
  `feat(api)!: drop the v1 chat endpoint`.
- Choose the type from the effect, not the file extension: tests only -> test,
  docs only -> docs, formatting or whitespace only -> style, dependency or
  tooling housekeeping -> chore or build, a revert -> revert, a mix -> feat
  or fix, whichever dominates.

LENGTH
- One line. Aim for 50 characters including the type and scope; never exceed
  72.
- No body, no footer, no bullet points, no numbered lists.

OUTPUT RULES
- Reply with the raw commit message and nothing else: no analysis, no
  explanation, no preamble, no quotes, no Markdown, no code fences, no
  "Commit message:" label.
- Use only facts visible in the diff (paths, identifiers, behaviour). Never
  invent a ticket number, issue reference or motivation.
- If the diff contains `# hunk truncated`, describe only the visible changes
  and do not speculate about the truncated parts.
- If the diff is empty or unreadable, reply with exactly:
  chore: update files
```

## Laravel sketch

```php
$validated = $request->validate([
    'model'        => 'nullable|string|max:100',
    'diff'         => 'required|string',
    'instructions' => 'nullable|string|max:20000',
]);

$system = $validated['instructions'] ?? config('gitmeh.instructions');

$response = Http::withToken(config('gitmeh.upstream_key'))
    ->timeout(60)
    ->post(config('gitmeh.upstream_url').'/chat/completions', [
        'model'       => ($validated['model'] ?? null) ?: config('gitmeh.upstream_model'),
        'messages'    => [
            ['role' => 'system', 'content' => $system],
            ['role' => 'user', 'content' => $validated['diff']],
        ],
        'temperature' => 0.3,
        'max_tokens'  => 4096,
    ]);

if ($response->failed()) {
    return response()->json(
        ['error' => ['message' => $response->json('error.message', 'upstream error')]],
        $response->status() >= 500 ? 502 : 422,
    );
}

$message = trim($response->json('choices.0.message.content', ''));

return response()->json([
    'choices' => [['message' => ['role' => 'assistant', 'content' => $message]]],
]);
```

Store the standard instructions in `config/gitmeh.php` as
`'instructions' => <<<'TXT' ... TXT`, not inline in the controller, and keep
`gitmeh-public-client` as the only accepted public token.
