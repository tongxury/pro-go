package volcengine

import (
	"context"
	"testing"
)

func TestClient_QueryMixCutTaskPreview(t *testing.T) {

	c := NewClient()

	c.QueryMixCutTaskPreview(context.Background(), QueryMixCutTaskPreviewParams{
		TaskKey: "9663242a-9adf-11f0-94f6-3436ac120034",
		//TaskKey: "5addab02-9ae7-11f0-a21d-3436ac120709",
	})
}
