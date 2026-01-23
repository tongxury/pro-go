package applestore

import (
	"context"
	"testing"
)

func TestClient_GetTransactionInfo(t *testing.T) {

	x := NewAppleStoreClient()

	x.GetTransactionInfo(context.Background(), "370002175250920")

}
