package arkr

import (
	"context"
	"fmt"
	"store/pkg/sdk/conv"
	"testing"
)

func TestClient_GetContentGenerationTask(t *testing.T) {

	c := NewClient()

	//task, err := c.GetContentGenerationTask(context.Background(), "cgt-20250926034712-64h4p")
	task, err := c.GetContentGenerationTask(context.Background(), "cgt-20250926035221-rlvsc")

	fmt.Println(conv.S2J(task), err)

}
