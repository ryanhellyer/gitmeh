//go:build integration

package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func skipWithoutGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not on PATH: %v", err)
	}
}

func mustGit(t *testing.T, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func mustWriteFile(t *testing.T, name, content string) {
	t.Helper()
	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func initWorktree(t *testing.T) {
	t.Helper()
	skipWithoutGit(t)

	root := t.TempDir()
	repo := filepath.Join(root, "work")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(repo)

	mustGit(t, "init", "-b", "main")
	mustGit(t, "config", "user.email", "integration-test@example.com")
	mustGit(t, "config", "user.name", "Integration Test")
}

// Checks StagedDiff() after a real staged edit: initial commit, change file, stage with git add, then expect a non-empty unified diff naming the file and showing the added line.
func TestIntegration_StagedDiff(t *testing.T) {
	initWorktree(t)

	mustWriteFile(t, "tracked.txt", "hello\n")
	mustGit(t, "add", "tracked.txt")
	mustGit(t, "commit", "-m", "init")

	mustWriteFile(t, "tracked.txt", "hello\nworld\n")
	mustGit(t, "add", "tracked.txt")

	diff, err := StagedDiff()
	if err != nil {
		t.Fatal(err)
	}
	if diff == "" {
		t.Fatal("expected non-empty staged diff")
	}
	if !strings.Contains(diff, "tracked.txt") {
		t.Fatalf("diff should name file; got:\n%s", diff)
	}
	if !strings.Contains(diff, "+world") {
		t.Fatalf("diff should contain added line; got:\n%s", diff)
	}
}

// Checks StagedDiff() by Staging an untracked file: one commit exists, a new file appears on disk only, AddAll runs, then StagedDiff must mention that file.
func TestIntegration_AddAll_stagesNewFile(t *testing.T) {
	initWorktree(t)

	mustWriteFile(t, "first.txt", "a\n")
	mustGit(t, "add", "first.txt")
	mustGit(t, "commit", "-m", "init")

	mustWriteFile(t, "second.txt", "new\n")
	if err := AddAll(); err != nil {
		t.Fatal(err)
	}

	diff, err := StagedDiff()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(diff, "second.txt") {
		t.Fatalf("staged diff should include new file; got:\n%s", diff)
	}
}

// Checks Commit with staged changes: after staging an edit, Commit records the given message; verify with git log -1.
func TestIntegration_Commit(t *testing.T) {
	initWorktree(t)

	mustWriteFile(t, "a.txt", "1\n")
	mustGit(t, "add", "a.txt")
	mustGit(t, "commit", "-m", "init")

	mustWriteFile(t, "a.txt", "2\n")
	mustGit(t, "add", "a.txt")

	if err := Commit("second message"); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command("git", "log", "-1", "--format=%s").CombinedOutput()
	if err != nil {
		t.Fatalf("git log: %v\n%s", err, out)
	}
	if got := strings.TrimSpace(string(out)); got != "second message" {
		t.Fatalf("HEAD subject: got %q want %q", got, "second message")
	}
}

// Exercises the full path with a real remote: bare repo as origin, initial commit pushed by hand, then staged edit + CommitAndPush; assert the latest commit subject on branch main in the bare repo (bare default HEAD stays on empty master, so we read main explicitly).
func TestIntegration_CommitAndPush(t *testing.T) {
	skipWithoutGit(t)

	root := t.TempDir()
	bare := filepath.Join(root, "origin.git")
	if err := os.MkdirAll(bare, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "init", "--bare", bare)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init --bare: %v\n%s", err, out)
	}

	repo := filepath.Join(root, "work")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(repo)

	mustGit(t, "init", "-b", "main")
	mustGit(t, "config", "user.email", "integration-test@example.com")
	mustGit(t, "config", "user.name", "Integration Test")
	mustGit(t, "remote", "add", "origin", bare)

	mustWriteFile(t, "file.txt", "v1\n")
	mustGit(t, "add", "file.txt")
	mustGit(t, "commit", "-m", "init")
	mustGit(t, "push", "origin", "HEAD")

	mustWriteFile(t, "file.txt", "v1\nv2\n")
	if err := AddAll(); err != nil {
		t.Fatal(err)
	}
	if err := CommitAndPush("pushed from integration test"); err != nil {
		t.Fatal(err)
	}

	// Bare repos often default HEAD to master (empty); first push created main.
	logOut, err := exec.Command("git", "-C", bare, "log", "-1", "main", "--format=%s").CombinedOutput()
	if err != nil {
		t.Fatalf("git log in bare: %v\n%s", err, logOut)
	}
	if got := strings.TrimSpace(string(logOut)); got != "pushed from integration test" {
		t.Fatalf("bare HEAD subject: got %q want %q", got, "pushed from integration test")
	}
}
