package trader

import (
	"store/pkg/sdk/chain/sol/jupiter"
	"store/pkg/sdk/chain/sol/quicknode"
	"store/pkg/sdk/chain/sol/solana/rpc"
	"store/pkg/sdk/chain/sol/solana/rpc/ws"
)

type Factory struct {
	solana              *rpc.Client
	jupiter             ITrader
	jupiterInstructions ITrader
	okx                 ITrader
	pumpfunInstructions ITrader
	raydium             ITrader
}

func NewFactory(solana *rpc.Client, solanaWS *ws.Client, qn *quicknode.Client) *Factory {

	jc := jupiter.NewSwapClient()
	//okc := okxapi.NewClient()

	return &Factory{
		solana:              solana,
		jupiter:             NewJupiterTrader(solana, jc),
		jupiterInstructions: NewJupiterInstructionsTrader(solana, jc),
		//okx:                 NewOKXTrader(solana, okc),
		pumpfunInstructions: NewPumpfunInstructionsTrader(solana, solanaWS, qn),
		raydium:             NewRaydiumTrader(solana),
	}
}

func (t *Factory) GetTrader() ITrader {
	return t.pumpfunInstructions
}

func (t *Factory) GetJupiterTrader() ITrader {
	return t.jupiter
}

func (t *Factory) GetRaydiumTrader() ITrader {
	return t.raydium
}

//func (t *Factory) GetPumpPortalTrader() ITrader {
//	return NewPumpPortalTrader(t.solana)
//}
