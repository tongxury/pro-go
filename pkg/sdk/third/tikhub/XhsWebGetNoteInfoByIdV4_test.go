package tikhub

import (
	"context"
	"testing"
)

func TestClient_XhsGetNoteByIdV4(t *testing.T) {

	c := NewClient()
	//
	c.XhsWebGetNoteInfoByIdV4(context.Background(), "5c6d9037000000001002ba15")
}
