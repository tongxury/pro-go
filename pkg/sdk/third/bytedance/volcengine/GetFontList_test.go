package volcengine

import (
	"context"
	"testing"
)

func TestClient_GetFontList(t *testing.T) {

	c := NewClient()

	c.GetFontList(
		context.Background(),
		GetFontListParams{},
	)
}
