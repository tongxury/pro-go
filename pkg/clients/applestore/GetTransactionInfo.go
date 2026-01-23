package applestore

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func (t *Client) GetTransactionInfo(ctx context.Context, transactionId string) {

	info, err := t.c.GetTransactionInfo(ctx, transactionId)
	if err != nil {
		return
	}

	var tokenClaims jwt.MapClaims
	err = t.s.ParseNotificationV2WithClaim(info.SignedTransactionInfo, &tokenClaims)
	if err != nil {
		return
	}

	fmt.Println(tokenClaims)
}
