package vikingdb

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_Delete(t *testing.T) {

	c := NewClient()

	ctx := context.Background()

	r := c.Delete(ctx, DeleteRequest{
		CollectionName: "template_commodity_coll",
		DeleteAll:      true,
	})

	fmt.Println(r)
}
