package sendandconfirmtransaction

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/chain/sol/solana/rpc"
	"store/pkg/sdk/chain/sol/solana/rpc/ws"
	"time"
)

func WaitForFinalized(
	ctx context.Context,
	wsClient *ws.Client,
	sig solana.Signature,
	timeout *time.Duration,
) (confirmed bool, err error) {
	sub, err := wsClient.SignatureSubscribe(
		sig,
		rpc.CommitmentFinalized,
	)

	log.Debugw("SignatureSubscribe ", time.Now().UnixMilli(), "signature", sig.String())

	if err != nil {
		return false, err
	}
	defer sub.Unsubscribe()

	if timeout == nil {
		t := 2 * time.Minute // random default timeout
		timeout = &t
	}

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(*timeout):
			return false, ErrTimeout
		case resp, ok := <-sub.Response():
			if !ok {
				return false, fmt.Errorf("subscription closed")
			}
			if resp.Value.Err != nil {
				// The transaction was confirmed, but it failed while executing (one of the instructions failed).
				return true, fmt.Errorf("confirmed transaction with execution error: %v", resp.Value.Err)
			} else {
				// Success! Confirmed! And there was no error while executing the transaction.
				return true, nil
			}
		case err := <-sub.Err():
			return false, err
		}
	}
}
