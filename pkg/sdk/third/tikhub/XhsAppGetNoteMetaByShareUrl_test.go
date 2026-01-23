package tikhub

import (
	"context"
	"testing"
)

func TestClient_XhsGetNoteByShareUrl(t *testing.T) {

	c := NewClient()
	//
	c.XhsGetNoteByShareUrl(context.Background(), "http://xhslink.com/a/hX42R5pYUGqbb")
}
