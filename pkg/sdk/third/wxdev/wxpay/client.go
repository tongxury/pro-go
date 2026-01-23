package wxpay

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	"time"
)

type Config struct {
	MchId    string
	SerialNo string
}
type Client struct {
	conf Config
}

func NewClient() *Client {
	return &Client{
		conf: Config{
			MchId:    "1725227465",
			SerialNo: "4989C9084A1DB67441E8E6CB6E201BED98EE3CEE",
		},
	}
}

type SignResult struct {
	Signature string
	NonceStr  string
	TimeStamp int64
}

func (t *Client) signMessage(message string) (string, error) {

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
-----END PRIVATE KEY-----
`)
	if err != nil {
		log.Fatal("load merchant private key error")
		return "", err
	}

	part2, err := utils.SignSHA256WithRSA(message, mchPrivateKey)
	if err != nil {
		return "", err
	}

	return part2, nil

}

func (t *Client) sign(method, url string, params string) (*SignResult, error) {
	nonceStr, _ := utils.GenerateNonce()
	ts := time.Now().Unix()

	part1 := fmt.Sprintf(`%s
%s
%d
%s
%s
`, method, url, ts, nonceStr, params)

	fmt.Println(part1)

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
-----END PRIVATE KEY-----
`)
	if err != nil {
		log.Fatal("load merchant private key error")
		return nil, err
	}

	hash := sha256.Sum256([]byte(part1))
	signature, err := rsa.SignPKCS1v15(nil, mchPrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}

	part21 := base64.StdEncoding.EncodeToString(signature)

	part2, err := utils.SignSHA256WithRSA(part1, mchPrivateKey)
	if err != nil {
		return nil, err
	}

	fmt.Println(part21)
	fmt.Println(part2)

	// 格式：WECHATPAY2-SHA256-RSA2048 mchid="商户号",nonce_str="随机串",signature="签名值",timestamp="时间戳",serial_no="证书序列号"
	auth := fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",signature="%s",timestamp="%d",serial_no="%s"`,
		t.conf.MchId,
		nonceStr,
		part2,
		ts,
		t.conf.SerialNo)

	fmt.Println(auth)
	return &SignResult{
		Signature: auth,
		NonceStr:  nonceStr,
		TimeStamp: ts,
	}, nil
}
