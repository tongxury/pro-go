package jupiter

import (
	"context"
	"fmt"
	"store/pkg/sdk/chain/sol/solana"
	"testing"
)

func TestGetTokenMetadata(t *testing.T) {

	c := NewJupiterClient()

	metadata, err := c.GetTokenMetadata(context.Background(), solana.SolMint.String())
	if err != nil {
		return
	}

	fmt.Println(metadata)
}
