package git

import (
	"errors"
	"os/exec"
	"strings"
)

// StagedDiff returns the unified diff of staged changes (git diff --cached).
func StagedDiff() (string, error) {
	out, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && len(exit.Stderr) > 0 {
			return "", errors.New(strings.TrimSpace(string(exit.Stderr)))
		}
		return "", err
	}
	return string(out), nil
}
