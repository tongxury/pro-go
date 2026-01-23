package tikhub

import (
	"context"
	"testing"
)

func TestClient_XhsGetNoteInfoById(t *testing.T) {

	c := NewClient()

	c.xhsAppFetchFeedNotesV2(context.Background(), "689254a20000000023023463")

}
