package arkr

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_CreateContentGenerationTask(t *testing.T) {

	c := NewClient()

	task, err := c.CreateContentGenerationTask(context.Background())

	fmt.Println(task, err)

}
