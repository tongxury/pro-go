package gemini

import "store/pkg/sdk"

type Config struct {
	BaseURL string
	Proxy   string
	Key     string
	Cache   sdk.ICache
}

type FactoryConfig struct {
	Configs []*Config
}
