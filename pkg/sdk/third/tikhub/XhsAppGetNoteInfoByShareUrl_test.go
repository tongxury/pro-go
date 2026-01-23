package tikhub

import (
	"context"
	"testing"
)

func TestClient_XhsAppGetNoteInfoShareUrl(t *testing.T) {

	c := NewClient()
	//
	c.XhsAppGetNoteInfoShareUrl(context.Background(), "https://xhslink.com/a/EZ4M9TwMA6c3")

}
