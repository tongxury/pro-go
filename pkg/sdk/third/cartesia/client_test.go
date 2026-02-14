package cartesia

import (
	"store/confs"
	"testing"
)

func getClient(t *testing.T) *Client {
	if confs.CartesiaKey == "" {
		t.Skip("Skipping test: CartesiaKey is empty")
	}
	return NewClient(confs.CartesiaKey)
}

func TestNewClient(t *testing.T) {
	client := NewClient("test-key")
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
}
