package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"os"
)

func main() {
	payload := `diff --git a/README.md b/README.md
--- a/README.md
+++ b/README.md
@@ -4,6 +4,11 @@
 ## Intro

 Short description.
+
+### GITMEH_PROBE_OK
+
+Document the /gitmeh smoke test: POST a unified diff as plain text.
+
 ## License
`

	gitPush()

	body, err := apiRequest(payload)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Println("Response:", body)
}

func apiRequest(payload string) (string, error) {
	url := "https://ai.hellyer.kiwi/gitmeh"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		// Fix: must return two values to match (string, error)
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

func gitPush() {

	// 1. git add --all
	err := runCommand("git", "add", "--all")
	if err != nil {
		fmt.Println("Error adding files:", err)
		return
	}

	// 2. git commit -m 'x'
	err = runCommand("git", "commit", "-m", "x")
	if err != nil {
		fmt.Println("Error committing:", err)
		return
	}

	// 3. git push origin master
	err = runCommand("git", "push", "origin", "master")
	if err != nil {
		fmt.Println("Error pushing:", err)
		return
	}

	fmt.Println("Git commands executed successfully!")
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	
	// This ensures you see the git output in your terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}
