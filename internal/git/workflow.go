package git

import "fmt"

// Publish stages all changes, commits, and pushes to origin/master.
func Publish() {
	// 1. git add --all
	if err := AddAll(); err != nil {
		fmt.Println("Error adding files:", err)
		return
	}

	// 2. git commit -m …
	if err := Commit("x"); err != nil {
		fmt.Println("Error committing:", err)
		return
	}

	// 3. git push origin master
	if err := PushOriginMaster(); err != nil {
		fmt.Println("Error pushing:", err)
		return
	}

	fmt.Println("Git commands executed successfully!")
}
