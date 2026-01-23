package tikhub

import (
	"context"
	"testing"
)

func TestClient_XhsGetNoteByIdV2(t *testing.T) {

	c := NewClient()
	//
	c.XhsWebGetNoteInfoByIdV2(context.Background(), "5c6d9037000000001002ba15")
}
