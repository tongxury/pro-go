package tikhub

import (
	"context"
	"testing"
)

func TestClient_XhsSearchNotes(t *testing.T) {

	c := NewClient()
	//
	c.xhsWebV2FetchSearchNotes(context.Background(), XhsSearchNotesParams{
		Keyword: "幽默油条果汁",
	})
}
