# Remaining Improvements

### #8 Diff size pre-flight check [LOW-MEDIUM]

No guard against sending massive diffs (documented ~618K token failure in README).

**Fix:** Check `git diff --cached` byte size before API call. Warn or refuse if beyond a configurable limit (e.g., 100KB).

### #9 Consistent error formatting [LOW]

- `fatalErr` prints bare `err.Error()` with no prefix
- `fatalMsg` prints bare message
- API errors use `"%s | %s"` format (HTTP status + body)
- No common `"gitmeh: error:"` prefix for grep-ability

**Fix:** Create a centralized error helper with consistent format.






# 10
| **No Makefile** | Build uses shell scripts instead of conventional `Makefile`. |
