package aiapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const httpTimeout = 60 * time.Second

// DefaultHTTPClient returns an [http.Client] configured with the timeout used
// for gitmeh's API requests.
func DefaultHTTPClient() *http.Client {
	return &http.Client{Timeout: httpTimeout}
}

// HTTPClientForChatBase returns a client with [httpTimeout]. If baseURL's
// hostname equals insecureHost, certificate verification is skipped to
// accommodate self-signed dev TLS (e.g. on a private server).
func HTTPClientForChatBase(baseURL string, insecureHost string) *http.Client {
	if !chatBaseSkipsTLSVerify(baseURL, insecureHost) {
		return DefaultHTTPClient()
	}
	tr, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return DefaultHTTPClient()
	}
	ct := tr.Clone()
	if ct.TLSClientConfig == nil {
		ct.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	} else if ct.TLSClientConfig.MinVersion == 0 {
		ct.TLSClientConfig.MinVersion = tls.VersionTLS12
	}
	ct.TLSClientConfig.InsecureSkipVerify = true //nolint:gosec // guarded by chatBaseSkipsTLSVerify
	return &http.Client{Timeout: httpTimeout, Transport: ct}
}

func chatBaseSkipsTLSVerify(baseURL string, host string) bool {
	u, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || u.Hostname() == "" {
		return false
	}
	return strings.EqualFold(u.Hostname(), host)
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
