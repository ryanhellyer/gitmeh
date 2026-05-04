package main

import (
	"bufio"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode/utf8"

	"gitmeh/internal/aiapi"
	"gitmeh/internal/config"
	"gitmeh/internal/git"
	"gitmeh/internal/version"

	"golang.org/x/term"
)

// helpText is filled at compile time from help.txt. The next line is a Go
// compiler directive (not ordinary documentation): it tells the toolchain to
// copy that file into the binary. import _ "embed" is required so the compiler
// enables //go:embed even though we do not reference embed.FS in code.
//
//go:embed help.txt
var helpText string

const commitMsgPrompt = "Commit message: "

func main() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-h", "--help":
			fmt.Print(helpText)
			return
		case "-v", "--version":
			fmt.Printf("gitmeh %s\n", version.Version)
			return
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := git.AddAll(); err != nil {
		fatalErr(err)
	}

	diff, err := git.StagedDiff()
	if err != nil {
		fatalErr(err)
	}
	if strings.TrimSpace(diff) == "" {
		fatalMsg("nothing staged to commit")
	}

	cfg := config.Load()
	if cfg.Chat.MaxDiffBytes > 0 && len(diff) > cfg.Chat.MaxDiffBytes {
		fatalMsg(fmt.Sprintf(
			"staged diff is %d bytes (max %d). Set GITMEH_MAX_DIFF_BYTES to increase the limit.",
			len(diff), cfg.Chat.MaxDiffBytes,
		))
	}
	httpClient := aiapi.HTTPClientForChatBase(cfg.Chat.BaseURL)
	msg, err := aiapi.CommitMessageOpenAIChat(ctx, httpClient, aiapi.OpenAIChatParams{
		BaseURL:        cfg.Chat.BaseURL,
		APIKey:         cfg.Chat.APIKey,
		Model:          cfg.Chat.Model,
		SystemPrompt:   cfg.Chat.Prompt,
		FallbackModels: cfg.Chat.FallbackModels,
	}, diff)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fatalMsg("cancelled")
		}
		fatalErr(err)
	}
	if strings.TrimSpace(msg) == "" {
		fatalMsg("empty commit message from API")
	}

	final, proceed, err := reviewCommitMessage(msg, os.Stdin, os.Stdout)
	if err != nil {
		fatalErr(err)
	}
	if !proceed {
		os.Exit(1)
	}

	if err := git.CommitAndPush(final); err != nil {
		fatalErr(err)
	}
	fmt.Println("Git commands executed successfully!")
}

// reviewCommitMessage loops until the user accepts (Y), aborts (n), or edits (e).
// After edit, pressing Enter on the edited line submits that message (ok true).
// ok is false when the user aborts without error.
// stdin and stdout are injected for tests (e.g. strings.NewReader, io.Discard).
// If stdout is nil, os.Stdout is used.
//nolint:errcheck // UI output writes — safe to ignore
func reviewCommitMessage(suggested string, stdin io.Reader, stdout io.Writer) (final string, ok bool, err error) {
	if stdout == nil {
		stdout = os.Stdout
	}

	current := strings.TrimSpace(suggested)
	rd := bufio.NewReader(stdin)

	for {
		fmt.Fprintln(stdout, "\nSuggested commit message:")
		fmt.Fprintln(stdout, current)
		fmt.Fprint(stdout, "\nAccept this message? [Y]es / [n]o / [e]dit: ")

		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return "", false, fmt.Errorf("no input (EOF)")
			}
			return "", false, err
		}

		ans := strings.ToLower(strings.TrimSpace(line))
		switch {
		case ans == "" || ans == "y" || ans == "yes":
			return current, true, nil
		case ans == "n" || ans == "no":
			fmt.Fprintln(stdout, "Aborted.")
			return "", false, nil
		case ans == "e" || ans == "edit":
			fmt.Fprintln(stdout)
			edited, err := readCommitMessageInline(current, stdin, rd, stdout)
			if err != nil {
				return "", false, err
			}
			edited = strings.TrimSpace(edited)
			if edited == "" {
				fmt.Fprintln(stdout, "Commit message is empty; keeping previous text.")
				continue
			}
			// Enter after editing submits this message; skip another Y/n/e round.
			return edited, true, nil
		default:
			fmt.Fprintln(stdout, "Please enter y, n, or e (or press Enter for yes).")
		}
	}
}

// stdinTerminalFD returns the FD to put in raw mode when stdin is an
// [os.File] connected to a terminal; otherwise ok is false.
func stdinTerminalFD(stdin io.Reader) (fd int, ok bool) {
	f, isFile := stdin.(*os.File)
	if !isFile {
		return 0, false
	}
	fd = int(f.Fd())
	if !term.IsTerminal(fd) {
		return 0, false
	}
	return fd, true
}

// readCommitMessageInline shows an editable single-line field prefilled with
// initial; Enter submits the current line. Ctrl+C aborts with an error.
// stdin must be the same reader rd wraps (e.g. bufio.NewReader(stdin)).
// stdout is used for prompts (inject io.Discard in tests). If nil, os.Stdout.
//nolint:errcheck // UI output writes — safe to ignore
func readCommitMessageInline(initial string, stdin io.Reader, rd *bufio.Reader, stdout io.Writer) (string, error) {
	if stdout == nil {
		stdout = os.Stdout
	}

	fd, useTTY := stdinTerminalFD(stdin)
	if !useTTY {
		fmt.Fprintf(stdout, "%s%s\n", commitMsgPrompt, initial)
		fmt.Fprint(stdout, "(not a terminal — press Enter to keep, or type a new message)\n> ")
		line, err := rd.ReadString('\n')
		if err != nil {
			return "", err
		}
		t := strings.TrimSpace(line)
		if t == "" {
			return initial, nil
		}
		return t, nil
	}

	old, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() { _ = term.Restore(fd, old) }()

	line := []rune(initial)
	pos := len(line)

	redraw := func() {
		left := ""
		if pos > 0 {
			left = string(line[:pos])
		}
		right := ""
		if pos < len(line) {
			right = string(line[pos:])
		}
		fmt.Fprintf(stdout, "\r\033[K%s%s%s", commitMsgPrompt, left, right)
		if n := len(right); n > 0 {
			fmt.Fprintf(stdout, "\033[%dD", n)
		}
	}

	redraw()

	for {
		r, size, err := rd.ReadRune()
		if err != nil {
			if err == io.EOF {
				fmt.Fprint(stdout, "\r\n")
				return "", io.EOF
			}
			return "", err
		}
		if r == utf8.RuneError && size == 1 {
			continue
		}

		switch r {
		case '\r', '\n':
			fmt.Fprint(stdout, "\r\n")
			return string(line), nil
		case 3: // Ctrl+C
			fmt.Fprint(stdout, "\r\n")
			return "", fmt.Errorf("interrupted")
		case 127, '\b':
			if pos > 0 {
				line = append(line[:pos-1], line[pos:]...)
				pos--
				redraw()
			}
		case 27: // ESC — arrow keys
			br, _, err := rd.ReadRune()
			if err != nil || br != '[' {
				continue
			}
			dir, _, err := rd.ReadRune()
			if err != nil {
				continue
			}
			switch dir {
			case 'D':
				if pos > 0 {
					pos--
					redraw()
				}
			case 'C':
				if pos < len(line) {
					pos++
					redraw()
				}
			}
		default:
			if r >= 32 {
				line = append(line[:pos], append([]rune{r}, line[pos:]...)...)
				pos++
				redraw()
			}
		}
	}
}
