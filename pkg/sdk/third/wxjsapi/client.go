package wxjsapi

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
)

type Client struct {
	c *core.Client
}

func NewClient() (*Client, error) {
	var mchID string = "1725227465"
	//var mchCertificateSerialNumber string = "1F523651CF31057ACACA4E5DEBD292F5081D6C15"
	var mchCertificateSerialNumber string = "4989C9084A1DB67441E8E6CB6E201BED98EE3CEE"
	//var pubKeyId = "PUB_KEY_ID_0117252274652025081700192211000800"
	var mchAPIv3Key string = "MIIEvQIBADANBgkqhkiG9w0BAQEFAASA"
	//var mchAPIv3Key string = "xxxxx"

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKey(`-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCfwWAGJ3KygYOL
J6Ar/YduVhblAAH0ygIa0KkPgMPPm4Ltaft/G83ChJ/99XHqqXhj0CWNnPEDHiYV
pmC9F97ulNXkDXmg7/7IyhRUvW8DC3V+p3zrpC1wWj5DcW5lR82sCuzM330Loh1y
HMYN3fG8qofpFqkChkrIEVHMvydxSjxHp72VfsKGQcv+x96JlMW7zdMur13bf8QV
okpKD4Qeo3crYPwQa6kvH2vg69VFkMqlUu5veANkEcH8GSFUS/8dfn/Ce3/SDzvn
LtPD50zGIlcEWLqdfI+Y9V5jzvyJ1wospZ6Sdpbb8kRF/TObPHYs98PlfUGcl8cn
QfCkDjGPAgMBAAECggEARO3VtbINnhGdpUv6oyBn7+Z7SGFhdrI2iEVhvIutcQVD
T5a76dMgS36X3aaqeTqX9FEZ8uk1YEgA3LSF5vDGwqA7TYO26AbFIvN2JD38CQ9H
qdcwYifbZy+4z7bAkDiT+FhaZedD9+IB3HunxNHvfJ0DuUGKuMkiuQIoDjSoGwvj
HChTJxOu8Le5g7hw5Wiq8HpAB4OgqRa0XKXpaqzNMn4qZ8nzJ1sywggX7oXWoDWA
UeHaBKXNagEiRLjcz2Qw/iB/wX8/LXR6nYRP7Aj8MEE6Es+v5nlCYssQRuokVvkV
NPOCsC1CLbI3GhPSc4uaWYB7KUCOt21KwVDZDiUKEQKBgQDTCoE62/ENFDFxfGhB
FB787lhLTXkZPSt8j1TKAD+901LVf15EpS9NgEuH2MQc1XWLgboltFE1opNH6JDB
UJlbkBd7L2ReF1HkOPJ/4WAbOG9qSlK21JqxyIZj+Ox/5jxvNMM0qwU8PpQ8g4p1
JLK7vEFaRDX9B5QouxYDVqjA+QKBgQDByeqJ78jFtfCAgH2l+GSNhTGpppFeqIdq
2O3/UZ9lGrtRBDtTxotK80pA/EnSp+RFLuTONnx5eK64KkOsgOCwg4eloflUzUG4
GSWLrVvjv1D3EK8NK+fhDgcueWJhZYx/qcsT2vfD13DMs8hWjhgcQpsMLN1YOlYV
ad3ZFh2wxwKBgG6cKeFcl3mgZM2zQ70gO6GbloFZSKg2zE0LnogFG2N1mAu4JwNZ
hHJdVLkNnrPyGRqWUqciXBH9dK4SsZPwl4BLBFOXIkbCeDRiuI7X5BRAPvz5mWKk
CbQ2gmFxfRsH5BLxF4LKRAMwVWdmFjKRmnAVGjeiWp2U1E2IyN/VEruZAoGAdJzZ
uwmE6pySTfGEKrSvZY3qFam7PpfxbTV++i4W2dNdNuJyBPar6X0/iJ2ImvAW6B7Q
5tpYywv5L6+XK54eF3n+zYgLrqEZU/wl4MiATCtbQGFUxXtPPNmhLrEyp2NhSY1W
O+t/PuVM5pGlE5jMH21hOdFhnO710Er1ieXKFg8CgYAh8RE3IZ9PO6bKSFzx/EI6
P4nZUO5chwIZMIHh/0wY3G44sUNX4I4y97GyIKaVwfPfVDMWAfGk4b4c/hia6NtY
AjeWSIuEE1T7uxYpwOLOEcxSWoo2KR+ZlVjHrQws2+mQwrrGT7stYR+DrQCw0rge
Y+QMEZbivRINJ4Rvf6V9PA==
-----END PRIVATE KEY-----`)
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	if mchPrivateKey.N.BitLen() != 2048 {
		log.Fatal("私钥长度不是2048位")
	}

	//	mchPublicKey, err := utils.LoadPublicKey(`
	//-----BEGIN PUBLIC KEY-----
	//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuNTiWkwSily2FwCP0EFL
	//t02Ww3YmzBZx69lp1A3NxAPSw0T8Zc1bRX1Qf453Zt88OHHbJMcfcdqlttgLS1Mj
	//2JODV5H1LadGZ3Xp9aMQBFY5rMNW4XgpBa/lqbM9YUPyAp1iQ1GlPsGb+Ckx/XRC
	//UZx+2xf+F7DL1dsyCc6s1JRnH3I4qPnm/TTYffaOSziPlLkpmVkafx8yWdZAe8a3
	//m5/zV2uDnaXg6fch4AG9qvmXDji9nmLWYGwTS1XLV0xFBej+luRpjjX4oLdXT62N
	//dz6jpkqbkJcdg/Qe74+VB2JnnuGCZY+FjXyaeXt6cpScFOKau4QRyZVU2gmznt3g
	//DQIDAQAB
	//-----END PUBLIC KEY-----
	//	`)
	//	if err != nil {
	//		log.Fatal("load merchant private key error")
	//	}

	ctx := context.Background()
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
		//option.WithWechatPayPublicKeyAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, pubKeyId, mchPublicKey),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		c: client,
	}, nil
}
