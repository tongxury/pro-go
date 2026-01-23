package okxapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"store/confs"
	"time"
)

type Config struct {
	Endpoint      string
	APIKey        string
	APISecret     string
	APIPassphrase string
}
type Client struct {
	config Config
}

func NewClient() *Client {
	return &Client{
		config: Config{
			Endpoint:      "https://www.okx.com",
			APIKey:        confs.OKXAPIKey,
			APISecret:     confs.OKXAPISecret,
			APIPassphrase: confs.OKXAPIPassphrase,
		},
	}
}

func (t *Client) authHeaders(method, url, body string) (map[string]string, error) {

	timestamp := time.Now().Format("2006-01-02T15:04:05.999Z")

	message := timestamp + method + url + body

	// 创建 HMAC SHA256 签名
	h := hmac.New(sha256.New, []byte(t.config.APISecret))
	h.Write([]byte(message))
	signature := h.Sum(nil)

	base64Signature := base64.StdEncoding.EncodeToString(signature)

	return map[string]string{
		"OK-ACCESS-KEY":        t.config.APIKey,
		"OK-ACCESS-TIMESTAMP":  timestamp,
		"OK-ACCESS-PASSPHRASE": t.config.APIPassphrase,
		"OK-ACCESS-SIGN":       base64Signature,
	}, nil
}
