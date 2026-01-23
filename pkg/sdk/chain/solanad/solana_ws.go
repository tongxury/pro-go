package solanad

import (
	"context"
	"store/pkg/confcenter"
	"store/pkg/sdk/chain/sol/solana/rpc/ws"
)

type WSClient struct {
	endpoint string
}

func (t *WSClient) Connect(ctx context.Context) (*ws.Client, error) {
	return ws.Connect(ctx, t.endpoint)
}

func NewWSClient(conf confcenter.SolanaConfig) *WSClient {

	//url := "wss://mainnet.helius-rpc.com/?api-key=57e03bd1-0ed0-44d2-8e88-6f0f71c82d70"
	//url = "wss://dark-radial-seed.solana-mainnet.quiknode.pro/95494d59ce7464c3b374dcae1d25c0a3cba837f2"
	//url = rpc.MainNetBeta_WS
	//url = "wss://solana-api.instantnodes.io/token-o1KC32zJcz7SwE9er6lL5zP7Ddy7jOFn"

	//client, err := ws.Connect(context.Background(), url)
	//if err != nil {
	//	client, err = ws.Connect(context.Background(), url)
	//}

	c := &WSClient{
		endpoint: conf.WSEndpoint,
	}

	return c
}
