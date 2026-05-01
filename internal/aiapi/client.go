package aiapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gitmeh/internal/config"
)

const httpTimeout = 20 * time.Second

var httpClient = &http.Client{Timeout: httpTimeout}

// commitMessageSpinner draws a simple ASCII spinner on stderr until stop is closed.
func commitMessageSpinner(stop <-chan struct{}, done chan<- struct{}) {
	defer close(done)

	frames := []string{"-", "\\", "|", "/"}
	ticker := time.NewTicker(90 * time.Millisecond)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-stop:
			_, _ = fmt.Fprint(os.Stderr, "\r\033[K")
			return
		case <-ticker.C:
			_, _ = fmt.Fprintf(os.Stderr, "\r\033[K%s Generating commit message...", frames[i%len(frames)])
			i++
		}
	}
}

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

	stopSpinner := make(chan struct{})
	spinnerDone := make(chan struct{})
	go commitMessageSpinner(stopSpinner, spinnerDone)
	defer func() {
		close(stopSpinner)
		<-spinnerDone
	}()

	resp, err := httpClient.Do(req)
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
