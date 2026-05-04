# Remaining Improvements

### #7 Graceful shutdown / context cancellation [LOW-MEDIUM]

- `CommitMessageOpenAIChat` uses `http.NewRequest` instead of `http.NewRequestWithContext`
- Spinner goroutine has no signal handling
- Ctrl+C during API call leaves terminal messy

**Fix:** Use `context.Context` throughout, add `os.Signal` handling in `main()`.

### #8 Diff size pre-flight check [LOW-MEDIUM]

No guard against sending massive diffs (documented ~618K token failure in README).

**Fix:** Check `git diff --cached` byte size before API call. Warn or refuse if beyond a configurable limit (e.g., 100KB).

### #9 Consistent error formatting [LOW]

- `fatalErr` prints bare `err.Error()` with no prefix
- `fatalMsg` prints bare message
- API errors use `"%s | %s"` format (HTTP status + body)
- No common `"gitmeh: error:"` prefix for grep-ability

**Fix:** Create a centralized error helper with consistent format.

### #10 Tests for TTY editing path [LOW]

`readCommitMessageInline` has complex terminal raw-mode logic (backspace, arrow keys, Ctrl+C) that is completely untested.

**Fix:** Extract editing logic into a testable function or use a pseudo-terminal library.

## Additional Observations

| Issue | Detail |
|---|---|
| **Go version** | `go 1.25.9` in go.mod is non-standard. Update to a stable release. |
| **Stale PLAN.md** | Describes improvements already implemented. Should be updated or removed. |
| **Unused Backend enum** | `Backend` has only one value (`BackendOpenAIChat`) — scaffolding for future backends. |
| **No .editorconfig** | Cross-editor consistency. |
| **No Makefile** | Build uses shell scripts instead of conventional `Makefile`. |
| **No commit message validation** | Prompt requests Conventional Commits format but no server-side validation. |
| **`verify-openai-chat.sh`** | Redundant Python dependency; could use `jq` or be replaced by a Go smoke test. |
