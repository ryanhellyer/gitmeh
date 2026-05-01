package main

import (
	"fmt"

	"gitmeh/internal/aiapi"
	"gitmeh/internal/git"
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

	if err := git.Publish(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Git commands executed successfully!")

	body, err := aiapi.Request(payload)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Println("Response:", body)
}
