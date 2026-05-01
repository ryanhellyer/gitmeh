package git

// Publish stages all changes, commits, and pushes to origin/master.
func Publish() error {
	if err := AddAll(); err != nil {
		return err
	}
	if err := Commit("x"); err != nil {
		return err
	}
	if err := PushOriginMaster(); err != nil {
		return err
	}
	return nil
}
