//go:build !integration

package aiapi

import (
	"testing"
)

func TestHTTPClientForChatBase_insecureTransportOnlyForDevHost(t *testing.T) {
	t.Parallel()

	dev := HTTPClientForChatBase("https://ai.hellyer.test/v1")
	if dev.Transport == nil {
		t.Fatal("expected custom transport for ai.hellyer.test")
	}
	prod := HTTPClientForChatBase("https://openrouter.ai/api/v1")
	if prod.Transport != nil {
		t.Fatal("expected default transport for non-dev host")
	}
}
