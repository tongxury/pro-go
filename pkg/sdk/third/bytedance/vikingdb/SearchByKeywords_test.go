package vikingdb

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_SearchByKeywords(t *testing.T) {

	c := NewClient()

	ctx := context.Background()

	r, _ := c.SearchByKeywords(ctx, SearchByKeywordsRequest{
		SearchRequest: SearchRequest{
			CollectionName: "segment_commodity_coll",
			IndexName:      "segment_commodity_idx",
		},
		Keywords: []string{"玉米"},
	})

	fmt.Println(r)
}
