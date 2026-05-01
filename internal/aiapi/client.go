package aiapi

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitmeh/internal/config"
)

// CommitMessage POSTs the unified diff as plain UTF-8 text and returns the
// response body as the commit message (leading/trailing whitespace trimmed).
// On non-2xx responses, the returned error includes the raw body as a quoted
// string for debugging the API.
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
	raw := string(bodyBytes)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("%s | raw body: %q", resp.Status, raw)
	}

	return strings.TrimSpace(raw), nil
}
