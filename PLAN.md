  Documentation and repo hygiene

  README vs implementation: The README still describes shell-era behavior (OpenRouter, curl, jq, env vars like GITMEH_DEFAULT_URL). The Go path only POSTs
   to the fixed URL in internal/config:

   internal/config/config.go lines 1-3

  package config
  const GitMehURL = "https://ai.hellyer.kiwi/gitmeh"

  and internal/aiapi/client.go always uses that constant. Treating the README as “what this repo actually ships today” would reduce confusion and support
  burden—either trim or clearly label historical sections.

  Noise in README: Top-of-file notes about API fixes / 502s read like a scratchpad; moving them to issues or an internal doc keeps the project entry point
   readable.

  ────────────────────────────────────────



  Testing strategy (there are currently no tests)

  High leverage, behavior-preserving additions:

  • internal/aiapi: Spin up httptest.Server that returns fixed bodies/status codes and assert trimming and error formatting for non-2xx. Use the same
    default timeout/client wiring so behavior stays identical.
  • reviewCommitMessage: Inject an io.Reader (or small interface) instead of hardcoding os.Stdin, then table-drive y/n/e/EOF paths. That does not change
     behavior when main still passes os.Stdin.
  • internal/git: Optional integration tests that set GIT_DIR/GIT_WORK_TREE to temp dirs and run real git—more setup, but catches regressions in
    runCommand / diff behavior.


  ────────────────────────────────────────



  CLI conventions (note the caveat)

  Today main prints errors and returns without setting a non-zero exit status. Many scripts expect failures to exit non-zero. Changing exit codes can
  affect automation, so it is not purely internal—but aligning with common CLI expectations is a typical follow-up once you decide that’s acceptable.

  ────────────────────────────────────────



  Engineering workflow

  CI: There is no .github/workflows (or similar) in the tree. A minimal pipeline—go test ./..., go vet, and optionally staticcheck / golangci-lint—catches
   regressions cheaply.

  Release/version: internal/version/version.go documents -ldflags overrides; pairing that with tagged releases and changelog discipline matches the doc.

  ────────────────────────────────────────



  Small implementation notes (optional polish)

  • run.go: Errors from exec sometimes become plain errors.New(stderr) without %w, so callers cannot use errors.Is/As on underlying exit errors.
    Wrapping with %w preserves messages while improving introspection—behavior for normal users stays the same.
  • Spinners and stderr: The spinner goroutine is fine; if you ever see races in tests, serializing spinner start/stop behind an interface helps.

  ────────────────────────────────────────



  Summary

  The codebase is already split sensibly (main, git, aiapi, config). The biggest wins without altering features are bringing README/help in line with the 
  Go rewrite, adding tests around HTTP and the prompt loop, injecting dependencies for testability, and optional CI. Exit codes are the one item that
  straddles “polish” vs “observable behavior,” so treat that as a deliberate product decision.
