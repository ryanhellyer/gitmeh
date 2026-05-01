package git

// CommitAndPush runs git commit with message then pushes the current branch to origin.
func CommitAndPush(message string) error {
	if err := Commit(message); err != nil {
		return err
	}
	if err := PushOriginHead(); err != nil {
		return err
	}
	return nil
}
