package git

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func runGit(args ...string) error {
	cmd := exec.Command("git", args...) //nolint:gosec

	cmd.Stdout = os.Stdout

	var stderr strings.Builder
	cmd.Stderr = &stderr

	err := cmd.Run()
	msg := strings.TrimSpace(stderr.String())

	if err != nil {
		if msg != "" {
			return errors.New(msg)
		}
		return err
	}

	return nil
}
