package tikhub

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_XhsGetUserNotes(t *testing.T) {

	c := NewClient()
	id, err := c.XhsGetUserNotes(context.Background(), "682daf8c000000000e02df67")
	if err != nil {
		return
	}

	fmt.Println(id)

}
