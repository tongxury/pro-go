package wxdev

import (
	"context"
	"testing"
)

func TestClient_GetOpenIdByCode(t *testing.T) {

	c := NewClient()

	c.GetOpenIdByCode(context.Background(), "0d3NJA1009mgNU1xQN000ijK6m1NJA15")

}
