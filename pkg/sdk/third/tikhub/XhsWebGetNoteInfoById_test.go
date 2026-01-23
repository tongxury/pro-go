package tikhub

import (
	"context"
	"testing"
)

func TestClient_XhsGetNoteById(t *testing.T) {

	c := NewClient()
	//
	c.XhsWebGetNoteInfoById(context.Background(), "5c6d9037000000001002ba15")
}
