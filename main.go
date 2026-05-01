package main

import (
	"fmt"
	"strings"

	"gitmeh/internal/aiapi"
	"gitmeh/internal/git"
)

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

	if err := git.CommitAndPush(msg); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Git commands executed successfully!")
}
