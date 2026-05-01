package git

import (
	"fmt"
	"os"
	"os/exec"
)

func Push() {
	// 1. git add --all
	err := runCommand("git", "add", "--all")
	if err != nil {
		fmt.Println("Error adding files:", err)
		return
	}

	// 2. git commit -m 'x'
	err = runCommand("git", "commit", "-m", "x")
	if err != nil {
		fmt.Println("Error committing:", err)
		return
	}

	// 3. git push origin master
	err = runCommand("git", "push", "origin", "master")
	if err != nil {
		fmt.Println("Error pushing:", err)
		return
	}

	fmt.Println("Git commands executed successfully!")
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	// This ensures you see the git output in your terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
