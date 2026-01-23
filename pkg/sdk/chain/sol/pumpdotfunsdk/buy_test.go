package pumpdotfunsdk_test

import (
	"context"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/chain/sol/solana/rpc"
	"store/pkg/sdk/chain/sol/solana/rpc/ws"
)

type TestConfig struct {
	rpcClient  *rpc.Client
	wsClient   *ws.Client
	PrivateKey solana.PrivateKey
	mint       solana.PublicKey
}

func GetTestConfig() TestConfig {
	testConfig := TestConfig{}
	//Endpoint:   "https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
	//	WSEndpoint: "wss://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
	testConfig.rpcClient = rpc.New("https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6")
	wsClient, err := ws.Connect(context.Background(), "wss://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6")
	if err != nil {
		panic(err)
	}
	testConfig.wsClient = wsClient
	testConfig.PrivateKey = solana.MustPrivateKeyFromBase58("BdMzr3Nh82SYbYVkKwToSPkkeMrR2KtSs2AKHvw7nGFqx4bgktokG2t2qwGQx1QZRM5azTGk2rAi6tyUk8Y1yyM")
	testConfig.mint = solana.MustPublicKeyFromBase58("7u9St4vhrEzn8yxE2WtfK3ip4hV2Z4omoLhaYvzypump")
	return testConfig
}

//func TestBuyToken(t *testing.T) {
//	testConfig := GetTestConfig()
//	//pumpdotfunsdk.SetDevnetMode()
//	sig, err := pumpdotfunsdk.CreateBuyInstructions(
//		testConfig.rpcClient,
//		//testConfig.wsClient,
//		testConfig.PrivateKey,
//		testConfig.mint,
//		10000,
//		10000,
//	)
//	if err != nil {
//		t.Fatalf("can't buy token: %s", err)
//	}
//	t.Logf("buy token signature: %s", sig)
//}
