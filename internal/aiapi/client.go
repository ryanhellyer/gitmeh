package aiapi

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"

	"gitmeh/internal/config"
)

func Request(payload string) (string, error) {
	req, err := http.NewRequest("POST", config.GitMehURL, bytes.NewBufferString(payload))
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

	return string(bodyBytes), nil
}
