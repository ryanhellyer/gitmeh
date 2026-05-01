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

// CommitMessage sends the unified diff as the request body and returns the
// suggested commit message from the API (plain text response).
func CommitMessage(diff string) (string, error) {
	req, err := http.NewRequest("POST", config.GitMehURL, bytes.NewBufferString(diff))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")

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

	return body, nil
}
