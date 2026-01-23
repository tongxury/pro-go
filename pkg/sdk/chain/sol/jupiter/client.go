package jupiter

type JupiterClient struct {
	orderEndpoint string
	apiEndpoint   string
	tokenEndpoint string
}

type SwapClient struct {
	endpoint string
}

func NewSwapClient() *SwapClient {
	return &SwapClient{
		//endpoint:      "https://quote-api.jup.ag/v6",
		//endpoint:      "http://157.90.94.185:8080",
		endpoint: "https://public.jupiterapi.com",
		//endpoint: "http://157.90.94.185:30384", // todo 不稳定
	}
}

// https://station.jup.ag/api-v6/get-quote
func NewJupiterClient() *JupiterClient {
	return &JupiterClient{
		orderEndpoint: "https://jup.ag",
		//apiEndpoint:   "https://api.jup.ag",
		apiEndpoint: "http://157.90.94.185:30384",
		//priceEndpoint: "https://price.jup.ag",
		tokenEndpoint: "https://tokens.jup.ag",
	}
}
