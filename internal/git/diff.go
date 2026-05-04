package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// StagedDiff returns the unified diff of staged changes (git diff --cached).
func StagedDiff() (string, error) {
	out, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && len(exit.Stderr) > 0 {
			return "", errors.New(strings.TrimSpace(string(exit.Stderr)))
		}
		return "", err
	}
	return string(out), nil
}

// StagedDiffTruncated returns the staged diff, truncated per-file if it
// exceeds maxBytes. If maxBytes <= 0 the full diff is returned unchanged.
func StagedDiffTruncated(maxBytes int) (string, error) {
	diff, err := StagedDiff()
	if err != nil {
		return "", err
	}
	if maxBytes <= 0 || len(diff) <= maxBytes {
		return diff, nil
	}
	return truncateByFile(diff, maxBytes), nil
}

type diffSection struct {
	header string // "diff --git ..." through "+++ b/..."
	hunks  string // everything after the header
}

func parseSections(diff string) []diffSection {
	if diff == "" {
		return nil
	}
	// Split on "\ndiff --git " to get file sections.
	parts := strings.Split(diff, "\ndiff --git ")
	sections := make([]diffSection, len(parts))
	for i, p := range parts {
		raw := p
		if i > 0 {
			raw = "diff --git " + p
		}
		sections[i] = splitHeader(raw)
	}
	return sections
}

// splitHeader splits a file section into its 4-line header and the hunk body.
func splitHeader(section string) diffSection {
	start := 0
	for i := 0; i < 4; i++ {
		off := strings.IndexByte(section[start:], '\n')
		if off < 0 {
			return diffSection{header: section}
		}
		start += off + 1
	}
	if start >= len(section) {
		return diffSection{header: section}
	}
	return diffSection{
		header: section[:start],
		hunks:  section[start:],
	}
}

func truncateByFile(diff string, maxBytes int) string {
	sections := parseSections(diff)
	if len(sections) <= 1 {
		// Single file — just truncate the hunk portion.
		// Header is small, keep it all. Truncate hunks to fit.
		h := sections[0].header
		avail := maxBytes - len(h)
		if avail <= 0 {
			return diff[:maxBytes] + "\n# diff truncated\n"
		}
		if len(sections[0].hunks) <= avail {
			return diff
		}
		return h + sections[0].hunks[:avail] + "\n# hunk truncated\n"
	}

	// Keep all headers — they're small and essential.
	var headerBuf strings.Builder
	for _, s := range sections {
		headerBuf.WriteString(s.header)
	}
	headerLen := headerBuf.Len()

	if headerLen >= maxBytes {
		// Can't fit all headers — skip some file sections entirely.
		var buf strings.Builder
		for _, s := range sections {
			if buf.Len()+len(s.header) > maxBytes {
				break
			}
			buf.WriteString(s.header)
		}
		return buf.String()
	}

	hunkBudget := maxBytes - headerLen

	// Calculate total hunk size across all files.
	totalHunkSize := 0
	for _, s := range sections {
		totalHunkSize += len(s.hunks)
	}
	if totalHunkSize == 0 {
		return headerBuf.String()
	}

	// Proportional allocation: each file gets a share of the hunk budget
	// based on its hunk size relative to the total.
	var buf strings.Builder
	buf.WriteString(headerBuf.String())

	remainingHunkBudget := hunkBudget
	remainingHunkSize := totalHunkSize

	for _, s := range sections {
		if len(s.hunks) == 0 {
			continue
		}
		alloc := int(int64(hunkBudget) * int64(len(s.hunks)) / int64(totalHunkSize))
		if alloc > remainingHunkBudget {
			alloc = remainingHunkBudget
		}
		if alloc < 20 {
			alloc = 20 // minimum meaningful hunk content
		}
		if alloc > len(s.hunks) {
			alloc = len(s.hunks)
		}
		if alloc >= len(s.hunks) {
			buf.WriteString(s.hunks)
		} else {
			buf.WriteString(s.hunks[:alloc])
			buf.WriteString("\n# hunk truncated\n")
		}
		remainingHunkBudget -= alloc
		remainingHunkSize -= len(s.hunks)
	}

	return buf.String()
}

func formatBytes(b int) string {
	if b < 1024 {
		return fmt.Sprintf("%d B", b)
	}
	if b < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
}
