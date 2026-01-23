package rpc

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"math"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/chain/sol/solana/programs/system"
	"time"
)

func (cl *Client) WaitForFinalized(sign solana.Signature, timeout time.Duration) (bool, error) {

	start := time.Now()
	for {

		time.Sleep(1000 * time.Millisecond)

		if time.Since(start) > timeout {
			return false, fmt.Errorf("timeout")
		}
		log.Debugw("GetSignatureStatuses ", "", "sign", sign.String())

		newCtx := context.Background()

		result, err := cl.GetSignatureStatuses(newCtx, false, sign)
		if err != nil {
			log.Errorw("GetSignatureStatuses err", err, "sign", sign.String())
			continue
		}
		if len(result.Value) > 0 && result.Value[0] != nil { // todo confirmed之后也有可能失败
			if result.Value[0].ConfirmationStatus == ConfirmationStatusConfirmed ||
				result.Value[0].ConfirmationStatus == ConfirmationStatusFinalized {
				return true, nil
			}
		}
	}

	return false, nil
}

type LatestBlock struct {
	Hash solana.Hash
	Slot uint64
}

type SendTransactionParams struct {
	Transaction *solana.Transaction
	PrivateKeys []string
	LatestBlock *LatestBlock
	NoWait      bool
}

//func (cl *Client) SendTransactionAndWaitUtilConfirmed(ctx context.Context, params SendTransactionParams) (string, bool, error) {
//
//	if params.Transaction.Message.RecentBlockhash.IsZero() {
//		if params.LatestBlock == nil {
//			latestBlockHashResult, err := cl.GetLatestBlockhash(ctx, CommitmentFinalized)
//			if err != nil {
//				return "", false, fmt.Errorf("could not get latest blockhash: %w", err)
//			}
//			//
//			//params.LatestBlock = &LatestBlock{
//			//	Hash: latestBlockHashResult.Value.Blockhash,
//			//	Slot: (latestBlockHashResult.Context.Slot),
//			//}
//
//			params.Transaction.Message.RecentBlockhash = latestBlockHashResult.Value.Blockhash
//		} else {
//			params.Transaction.Message.RecentBlockhash = params.LatestBlock.Hash
//		}
//
//	}
//
//	err := params.Transaction.Signing(params.PrivateKeys...)
//	if err != nil {
//		return "", false, err
//	}
//
//	var maxRetries uint = 3
//	sig, err := cl.SendTransactionWithOpts(ctx, params.Transaction, TransactionOpts{
//		//Encoding:            solana.EncodingBase64Zstd,
//		MaxRetries: &maxRetries,
//		//MinContextSlot:      &params.LatestBlock.Slot,
//		PreflightCommitment: CommitmentProcessed,
//	})
//	if err != nil {
//		return "", false, fmt.Errorf("could not send transaction: %w", err)
//	}
//
//	sign := sig.String()
//
//	if !params.NoWait {
//
//		time.Sleep(1 * time.Second)
//		for i := 0; i < 10; i++ {
//			time.Sleep(1 * time.Second)
//			log.Debugw("GetSignatureStatuses ", i)
//
//			newCtx := context.Background()
//
//			result, err := cl.GetSignatureStatuses(newCtx, false, sig)
//			if err != nil {
//				log.Errorw("GetSignatureStatuses err", err, "sign", sig.String())
//				continue
//			}
//			if len(result.Value) > 0 && result.Value[0] != nil { // todo confirmed之后也有可能失败
//				if result.Value[0].ConfirmationStatus == ConfirmationStatusConfirmed ||
//					result.Value[0].ConfirmationStatus == ConfirmationStatusFinalized {
//					return sign, true, nil
//				}
//			}
//		}
//
//	}
//
//	return sign, false, nil
//}

func (cl *Client) CreateTransferTransaction(ctx context.Context, from, to string, amount float64) (*solana.Transaction, error) {

	fromPubKey := solana.MustPublicKeyFromBase58(from)

	blockhash, err := cl.GetLatestBlockhash(ctx, CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	am := uint64(amount * math.Pow10(9))

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				am,
				fromPubKey,
				solana.MustPublicKeyFromBase58(to),
			).Build(),
		},
		blockhash.Value.Blockhash,
		solana.TransactionPayer(fromPubKey),
	)

	return tx, nil
}
