package git

// PushOriginMaster runs `git push origin master`.
func PushOriginMaster() error {
	return runCommand("git", "push", "origin", "master")
}
