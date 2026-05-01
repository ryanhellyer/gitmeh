package git

// CommitAndPush runs git commit with message then pushes to origin/master.
func CommitAndPush(message string) error {
	if err := Commit(message); err != nil {
		return err
	}
	if err := PushOriginMaster(); err != nil {
		return err
	}
	return nil
}
