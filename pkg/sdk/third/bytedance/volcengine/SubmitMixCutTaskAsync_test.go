package volcengine

import (
	"context"
	"testing"
)

func TestClient_SubmitMixCutTaskAsync(t *testing.T) {

	c := NewClient()

	c.SubmitMixCutTaskAsync(context.Background(), SubmitMixCutTaskAsyncParams{
		//TaskKey:  "9663242a-9adf-11f0-94f6-3436ac120034",
		TaskKey:  "5addab02-9ae7-11f0-a21d-3436ac120709",
		GroupIds: []int{0},
	})
}
