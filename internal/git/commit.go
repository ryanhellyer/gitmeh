package git

func Commit(message string) error {
	return runGit("commit", "-m", message)
}
