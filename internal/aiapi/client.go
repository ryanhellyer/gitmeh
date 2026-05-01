package aiapi

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitmeh/internal/config"
)

// CommitMessage POSTs the unified diff as plain UTF-8 text and returns the
// response body as the commit message. The API must respond with plain text
// only (same idea as curl: one line like "Add … and document …", not JSON or
// Markdown).
func CommitMessage(diff string) (string, error) {
	req, err := http.NewRequest("POST", config.GitMehURL, bytes.NewBufferString(diff))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
	req.Header.Set("Accept", "text/plain")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	body := strings.TrimSpace(string(bodyBytes))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if body != "" {
			return "", fmt.Errorf("%s: %s", resp.Status, body)
		}
		return "", fmt.Errorf("%s", resp.Status)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "" {
		if strings.Contains(strings.ToLower(ct), "application/json") {
			return "", errors.New("API Content-Type is JSON; expected plain text")
		}
	}

	if strings.HasPrefix(body, "```") {
		return "", errors.New("API returned markdown fences; expected plain text only")
	}

	return body, nil
}
