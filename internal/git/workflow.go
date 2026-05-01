package git

import "fmt"

// Publish stages all changes, commits, and pushes to origin/master.
func Publish() error {
	// 1. git add --all
	if err := AddAll(); err != nil {
		return fmt.Errorf("git add: %w", err)
	}

	// 2. git commit -m …
	if err := Commit("x"); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}

	// 3. git push origin master
	if err := PushOriginMaster(); err != nil {
		return fmt.Errorf("git push: %w", err)
	}

	return nil
}
