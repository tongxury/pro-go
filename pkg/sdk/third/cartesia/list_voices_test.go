package cartesia

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_ListVoices(t *testing.T) {
	client := getClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ListVoices(ctx, &ListVoicesRequest{Language: "zh"})
	if err != nil {
		t.Fatalf("ListVoices failed: %v", err)
	}

	assert.NotEmpty(t, resp.Data, "Expected voices list to be non-empty")

	for _, v := range resp.Data {
		assert.NotEmpty(t, v.Id, "Voice ID should not be empty")
		assert.NotEmpty(t, v.Name, "Voice Name should not be empty")
		t.Logf("Found voice: %s (%s) - Language: %s", v.Name, v.Id, v.Language)
	}
}
