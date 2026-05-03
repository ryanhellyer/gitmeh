package aiapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const httpTimeout = 20 * time.Second

// DefaultHTTPClient returns an [http.Client] configured with the timeout used
// for gitmeh's API requests. Pass it with your endpoint URL to [CommitMessage],
// or supply your own client (for example in tests).
func DefaultHTTPClient() *http.Client {
	return &http.Client{Timeout: httpTimeout}
}

// stderrCommitSpinner draws a simple ASCII spinner on stderr until stop is closed.
func stderrCommitSpinner(stop <-chan struct{}, done chan<- struct{}) {
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

// withGeneratingCommitSpinner runs fn while showing a stderr spinner until fn returns.
func withGeneratingCommitSpinner(fn func() (string, error)) (string, error) {
	stopSpinner := make(chan struct{})
	spinnerDone := make(chan struct{})
	go stderrCommitSpinner(stopSpinner, spinnerDone)
	defer func() {
		close(stopSpinner)
		<-spinnerDone
	}()
	return fn()
}

// CommitMessage POSTs the unified diff to endpoint as plain UTF-8 text and
// returns the response body as the commit message (leading/trailing
// whitespace trimmed). On non-2xx responses, the returned error includes the
// raw body as a quoted string for debugging the API. client must not be nil.
func CommitMessage(client *http.Client, endpoint, diff string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("http client is nil")
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(diff))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
	req.Header.Set("Accept", "text/plain")

	return withGeneratingCommitSpinner(func() (string, error) {
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
	})
}
