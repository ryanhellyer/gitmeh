package git

func Commit(message string) error {
	return runCommand("git", "commit", "-m", message)
}
