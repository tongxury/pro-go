package applestore

import (
	"github.com/awa/go-iap/appstore"
	"github.com/awa/go-iap/appstore/api"
)

type Client struct {
	c *api.StoreClient
	s *appstore.Client
}

func NewAppleStoreClient() *Client {

	return &Client{
		s: appstore.New(),
		c: api.NewStoreClient(&api.StoreConfig{
			KeyContent: []byte(`
-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgZA5NJ8U9yV/AUA68
r/DnYDzj7Z4LE5jwtcjOB6yUtvmgCgYIKoZIzj0DAQehRANCAATv8pWKckBS7xOZ
UMythD/b6S/OQBJUxBuVv6OmmALFJriJptOdkTrUABnMqqDXWhWK0fDB7oHYTz0b
HyVyFRLd
-----END PRIVATE KEY-----
`),
			KeyID:              "B339XC9FB7",
			BundleID:           "com.tuturduck.veogoapp",
			Issuer:             "69a6de94-244d-47e3-e053-5b8c7c11a4d1",
			Sandbox:            false,
			TokenIssuedAtFunc:  nil,
			TokenExpiredAtFunc: nil,
			HostDebug:          "",
		}),
	}
}
