package fiware

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(Host("https://example.de/"))
	if client.Host != "https://example.de" {
		t.Errorf("host is not correct")
	}
}
