//go:build !integration

package aiapi

import (
	"testing"
)

func TestHTTPClientForChatBase_insecureTransportOnlyForDevHost(t *testing.T) {
	t.Parallel()

	dev := HTTPClientForChatBase("https://example.test/v1", "example.test")
	if dev.Transport == nil {
		t.Fatal("expected custom transport for host with matching insecure flag")
	}
	prod := HTTPClientForChatBase("https://openrouter.ai/api/v1", "example.test")
	if prod.Transport != nil {
		t.Fatal("expected default transport for unmatched host")
	}
}
