package git

func AddAll() error {
	return runGit("add", "--all")
}
