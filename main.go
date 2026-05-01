package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"gitmeh/internal/aiapi"
	"gitmeh/internal/git"

	"golang.org/x/term"
)

const commitMsgPrompt = "Commit message: "

func main() {
	if err := git.AddAll(); err != nil {
		fmt.Println(err)
		return
	}

	diff, err := git.StagedDiff()
	if err != nil {
		fmt.Println(err)
		return
	}
	if strings.TrimSpace(diff) == "" {
		fmt.Println("nothing staged to commit")
		return
	}

	msg, err := aiapi.CommitMessage(diff)
	if err != nil {
		fmt.Println(err)
		return
	}
	if strings.TrimSpace(msg) == "" {
		fmt.Println("empty commit message from API")
		return
	}

	final, proceed, err := reviewCommitMessage(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !proceed {
		return
	}

	if err := git.CommitAndPush(final); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Git commands executed successfully!")
}

// reviewCommitMessage loops until the user accepts (Y), aborts (n), or edits (e)
// and then reviews again. ok is false when the user aborts without error.
func reviewCommitMessage(suggested string) (final string, ok bool, err error) {
	current := strings.TrimSpace(suggested)
	rd := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nSuggested commit message:")
		fmt.Println(current)
		fmt.Print("\nAccept this message? [Y]es / [n]o / [e]dit: ")

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
			fmt.Println("Aborted.")
			return "", false, nil
		case ans == "e" || ans == "edit":
			fmt.Println()
			edited, err := readCommitMessageInline(current, rd)
			if err != nil {
				return "", false, err
			}
			edited = strings.TrimSpace(edited)
			if edited == "" {
				fmt.Println("Commit message is empty; keeping previous text.")
				continue
			}
			current = edited
			continue
		default:
			fmt.Println("Please enter y, n, or e (or press Enter for yes).")
		}
	}
}

// readCommitMessageInline shows an editable single-line field prefilled with
// initial; Enter submits the current line. Ctrl+C aborts with an error.
// rd must be the same bufio.Reader used for the surrounding menu reads.
func readCommitMessageInline(initial string, rd *bufio.Reader) (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		fmt.Printf("%s%s\n", commitMsgPrompt, initial)
		fmt.Print("(not a terminal — press Enter to keep, or type a new message)\n> ")
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
		fmt.Printf("\r\033[K%s%s%s", commitMsgPrompt, left, right)
		if n := len(right); n > 0 {
			fmt.Printf("\033[%dD", n)
		}
	}

	redraw()

	for {
		r, size, err := rd.ReadRune()
		if err != nil {
			if err == io.EOF {
				fmt.Print("\r\n")
				return "", io.EOF
			}
			return "", err
		}
		if r == utf8.RuneError && size == 1 {
			continue
		}

		switch r {
		case '\r', '\n':
			fmt.Print("\r\n")
			return string(line), nil
		case 3: // Ctrl+C
			fmt.Print("\r\n")
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
