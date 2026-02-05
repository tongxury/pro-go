package gemini

import "store/pkg/sdk"

type Config struct {
	BaseURL         string
	Proxy           string
	Key             string
	Project         string
	Location        string
	APIVersion      string
	CredentialsJSON string
	Cache           sdk.ICache
}

type FactoryConfig struct {
	Configs []*Config
}
