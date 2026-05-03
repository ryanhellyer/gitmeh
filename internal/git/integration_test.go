//go:build integration

package git

import "testing"

// TestIntegrationTagPlaceholder exists so `go test -tags=integration ./...`
// exercises the integration build; it does not call git.
func TestIntegrationTagPlaceholder(t *testing.T) {
	t.Log("integration build tag enabled")
}
