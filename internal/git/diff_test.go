//go:build !integration

package git

import (
	"strings"
	"testing"
)

func TestSplitHeader(t *testing.T) {
	t.Parallel()

	section := "diff --git a/a.go b/a.go\nindex abc..def 100644\n--- a/a.go\n+++ b/a.go\n@@ -1 +1 @@\n-foo\n+bar\n"

	got := splitHeader(section)
	if got.header != "diff --git a/a.go b/a.go\nindex abc..def 100644\n--- a/a.go\n+++ b/a.go\n" {
		t.Errorf("header: %q", got.header)
	}
	if got.hunks != "@@ -1 +1 @@\n-foo\n+bar\n" {
		t.Errorf("hunks: %q", got.hunks)
	}
}

func TestSplitHeader_noHunks(t *testing.T) {
	t.Parallel()

	section := "diff --git a/a.go b/a.go\nindex abc..def 100644\n--- a/a.go\n+++ b/a.go\n"

	got := splitHeader(section)
	if got.header != section {
		t.Errorf("expected full section as header: %q", got.header)
	}
}

func TestParseSections(t *testing.T) {
	t.Parallel()

	diff := "diff --git a/a.go b/a.go\nindex a..b 100644\n--- a/a.go\n+++ b/a.go\n@@ -1 +1 @@\n-foo\n+bar\n" +
		"diff --git b/b.go b/b.go\nindex c..d 100644\n--- b/b.go\n+++ b/b.go\n@@ -2 +2 @@\n-baz\n+qux\n"

	sections := parseSections(diff)
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}
	if !strings.Contains(sections[0].header, "a.go") {
		t.Errorf("section 0 header: %q", sections[0].header)
	}
	if !strings.Contains(sections[1].header, "b.go") {
		t.Errorf("section 1 header: %q", sections[1].header)
	}
}

func TestParseSections_single(t *testing.T) {
	t.Parallel()

	diff := "diff --git a/a.go b/a.go\nindex a..b 100644\n--- a/a.go\n+++ b/a.go\n@@ -1 +1 @@\n-foo\n+bar\n"

	sections := parseSections(diff)
	if len(sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(sections))
	}
	if !strings.Contains(sections[0].header, "a.go") {
		t.Errorf("header: %q", sections[0].header)
	}
	if sections[0].hunks != "@@ -1 +1 @@\n-foo\n+bar\n" {
		t.Errorf("hunks: %q", sections[0].hunks)
	}
}

func TestParseSections_empty(t *testing.T) {
	t.Parallel()

	sections := parseSections("")
	if len(sections) != 0 {
		t.Fatalf("expected 0 sections, got %d", len(sections))
	}
}

func TestTruncateByFile_singleFits(t *testing.T) {
	t.Parallel()

	diff := "diff --git a/a.go b/a.go\nindex a..b 100644\n--- a/a.go\n+++ b/a.go\n@@ -1 +1 @@\n-foo\n+bar\n"

	got := truncateByFile(diff, len(diff)+100)
	if got != diff {
		t.Errorf("expected full diff, got %q", got)
	}
}

func TestTruncateByFile_singleTooBig(t *testing.T) {
	t.Parallel()

	diff := "diff --git a/a.go b/a.go\nindex a..b 100644\n--- a/a.go\n+++ b/a.go\n@@ -1 +1 @@\n-foo\n+bar\n"

	got := truncateByFile(diff, 80)
	if !strings.Contains(got, "truncated") {
		t.Errorf("expected truncated: %q", got)
	}
	if !strings.Contains(got, "diff --git") {
		t.Errorf("expected header: %q", got)
	}
}

func TestTruncateByFile_multiFiles(t *testing.T) {
	t.Parallel()

	diff := "diff --git a/a.go b/a.go\nindex a..b 100644\n--- a/a.go\n+++ b/a.go\n@@ -1 +1 @@\n-foo\n+bar\n" +
		"diff --git b/b.go b/b.go\nindex c..d 100644\n--- b/b.go\n+++ b/b.go\n@@ -2 +2 @@\n-baz\n+qux\n"

	// 140 bytes: fits both headers (~132 B) with a little hunk room
	got := truncateByFile(diff, 140)
	if !strings.Contains(got, "a.go") {
		t.Errorf("missing a.go: %q", got)
	}
	if !strings.Contains(got, "b.go") {
		t.Errorf("missing b.go: %q", got)
	}
	// Should have truncated markers since both hunks don't fully fit
	if !strings.Contains(got, "truncated") {
		t.Errorf("expected truncation: %q", got)
	}
}

func TestTruncateByFile_notTooBig(t *testing.T) {
	t.Parallel()

	diff := "diff --git a/a.go b/a.go\nindex a..b 100644\n--- a/a.go\n+++ b/a.go\n@@ -1 +1 @@\n-foo\n+bar\n"

	got := truncateByFile(diff, 200)
	if got != diff {
		t.Errorf("expected full diff, got %q", got)
	}
}
