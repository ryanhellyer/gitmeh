package git

// PushOriginHead runs `git push origin HEAD`, pushing the current branch to
// origin using the same branch name on the remote.
func PushOriginHead() error {
	return runCommand("git", "push", "origin", "HEAD")
}
