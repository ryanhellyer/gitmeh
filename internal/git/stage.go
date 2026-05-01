package git

func AddAll() error {
	return runCommand("git", "add", "--all")
}
