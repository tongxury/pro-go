package volcengine

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_CreateUrlMaterial(t *testing.T) {

	c := NewClient()

	r, _ := c.CreateUrlMaterial(
		context.Background(),
		CreateMaterialParams{
			MediaFirstCategory: "video",
			Title:              "videoTitle22",
			MaterialUrl:        "https://yoozyres.tos-cn-shanghai.volces.com/5002300f0b4b343ffbadfe5a5fff7483.mp4",
			Wait:               true,
		},
	)

	fmt.Println(r)
}
