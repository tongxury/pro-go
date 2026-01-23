package pumpdotfunsdk

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"

	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/chain/sol/solana/rpc"
)

// BondingCurveData holds the relevant information decoded from the on-chain data.
type BondingCurveData struct {
	RealTokenReserves    *big.Int
	VirtualTokenReserves *big.Int
	VirtualSolReserves   *big.Int
}

func (b *BondingCurveData) String() string {
	return fmt.Sprintf("RealTokenReserves=%s, VirtualTokenReserves=%s, VirtualSolReserves=%s", b.RealTokenReserves, b.VirtualTokenReserves, b.VirtualSolReserves)
}

// fetchBondingCurve fetches the bonding curve data from the blockchain and decodes it.
func fetchBondingCurve(rpcClient *rpc.Client, bondingCurvePubKey solana.PublicKey) (*BondingCurveData, error) {
	accountInfo, err := rpcClient.GetAccountInfoWithOpts(context.TODO(), bondingCurvePubKey, &rpc.GetAccountInfoOpts{Encoding: solana.EncodingBase64, Commitment: rpc.CommitmentProcessed})
	if err != nil || accountInfo.Value == nil {
		return nil, fmt.Errorf("FBCD: failed to get account info: %w", err)
	}

	data := accountInfo.Value.Data.GetBinary()
	if len(data) < 24 {
		return nil, fmt.Errorf("FBCD: insufficient data length")
	}

	// Decode the bonding curve data assuming it follows little-endian format
	realTokenReserves := big.NewInt(0).SetUint64(binary.LittleEndian.Uint64(data[0:8]))
	virtualTokenReserves := big.NewInt(0).SetUint64(binary.LittleEndian.Uint64(data[8:16]))
	virtualSolReserves := big.NewInt(0).SetUint64(binary.LittleEndian.Uint64(data[16:24]))

	return &BondingCurveData{
		RealTokenReserves:    realTokenReserves,
		VirtualTokenReserves: virtualTokenReserves,
		VirtualSolReserves:   virtualSolReserves,
	}, nil
}
