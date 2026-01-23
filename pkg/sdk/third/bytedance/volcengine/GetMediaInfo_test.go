package volcengine

import (
	"context"
	"fmt"
	"store/pkg/sdk/conv"
	"testing"
)

func TestClient_GetMediaInfo(t *testing.T) {

	c := NewClient()

	info, err := c.GetMediaInfo(
		context.Background(),
		GetMediaInfoParams{
			MediaIds:  []string{"7554715494201311278"},
			MediaType: 1,
		},
	)

	fmt.Println(conv.S2J(info), err)
}
