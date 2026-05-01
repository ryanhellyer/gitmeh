package git

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stdout = os.Stdout

	var stderr strings.Builder
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	err := cmd.Run()
	if err == nil {
		return nil
	}

	msg := strings.TrimSpace(stderr.String())
	if msg != "" {
		return fmt.Errorf("%w: %s", err, msg)
	}
	return err
}
