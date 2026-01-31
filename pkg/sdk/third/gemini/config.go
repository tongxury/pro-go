package gemini

import "store/pkg/sdk"

type Config struct {
	BaseURL  string
	Proxy    string
	Key      string
	Project  string // GCP Project ID (required for Vertex AI)
	Location string // GCP Location/Region (required for Vertex AI, e.g. "us-central1")
	Cache    sdk.ICache
}

type FactoryConfig struct {
	Configs []*Config
}
