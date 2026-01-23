package volcengine

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"store/confs"
)

type Client struct {
	endpoint string
}

const (

	// AccessKeyID 请求凭证，从访问控制申请
	AccessKeyID     = confs.BytedanceAccessKeyID
	SecretAccessKey = confs.BytedanceSecretAccessKey

	// 请求地址
	Addr = "https://icp.volcengineapi.com"
	Path = "/" // 路径，不包含 Query

	// 请求接口信息
	//Service = "iccloud_muse"
	Region = "cn-north"
	//Action  = "SearchTemplate"
	Version = "2025-03-26"
)

func NewClient() *Client {
	return &Client{
		endpoint: "https://icp.volcengineapi.com",
	}
}

// https://bytedance.larkoffice.com/wiki/PeZGwrAWoiRJ5bkIwvecw9FHnl2
func (t *Client) hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func (t *Client) getSignedKey(secretKey, date, region, service string) []byte {
	kDate := t.hmacSHA256([]byte(secretKey), date)
	kRegion := t.hmacSHA256(kDate, region)
	kService := t.hmacSHA256(kRegion, service)
	kSigning := t.hmacSHA256(kService, "request")

	return kSigning
}

func (t *Client) hashSHA256(data []byte) []byte {
	hash := sha256.New()
	if _, err := hash.Write(data); err != nil {
		//log.Printf("input hash err:%s", err.Error())
	}

	return hash.Sum(nil)
}

type Req struct {
	Version string
	Service string
	Action  string
	Method  string
	Params  map[string]string
	Body    []byte
}

func (t *Client) doRequest(ctx context.Context, req Req) ([]byte, error) {

	queries := url.Values{}
	for k, v := range req.Params {
		queries.Set(k, fmt.Sprintf("%v", v))
	}

	if req.Version == "" {
		req.Version = Version
	}

	if req.Service == "" {
		req.Service = "iccloud_muse"
	}

	// 1. 构建请求
	queries.Set("Action", req.Action)
	queries.Set("Version", req.Version)
	requestAddr := fmt.Sprintf("%s?%s", Addr, queries.Encode())

	request, err := http.NewRequest(req.Method, requestAddr, bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, fmt.Errorf("bad request: %w", err)
	}

	// 2. 构建签名材料
	now := time.Now()
	date := now.UTC().Format("20060102T150405Z")
	authDate := date[:8]
	request.Header.Set("X-Date", date)

	payload := hex.EncodeToString(t.hashSHA256(req.Body))
	request.Header.Set("X-Content-Sha256", payload)
	request.Header.Set("Content-Type", "application/json")

	queryString := strings.Replace(queries.Encode(), "+", "%20", -1)
	signedHeaders := []string{"host", "x-date", "x-content-sha256", "content-type"}
	var headerList []string
	for _, header := range signedHeaders {
		if header == "host" {
			headerList = append(headerList, header+":"+request.Host)
		} else {
			v := request.Header.Get(header)
			headerList = append(headerList, header+":"+strings.TrimSpace(v))
		}
	}
	headerString := strings.Join(headerList, "\n")

	canonicalString := strings.Join([]string{
		req.Method,
		Path,
		queryString,
		headerString + "\n",
		strings.Join(signedHeaders, ";"),
		payload,
	}, "\n")

	hashedCanonicalString := hex.EncodeToString(t.hashSHA256([]byte(canonicalString)))

	credentialScope := authDate + "/" + Region + "/" + req.Service + "/request"
	signString := strings.Join([]string{
		"HMAC-SHA256",
		date,
		credentialScope,
		hashedCanonicalString,
	}, "\n")

	// 3. 构建认证请求头
	signedKey := t.getSignedKey(SecretAccessKey, authDate, Region, req.Service)
	signature := hex.EncodeToString(t.hmacSHA256(signedKey, signString))

	authorization := "HMAC-SHA256" +
		" Credential=" + AccessKeyID + "/" + credentialScope +
		", SignedHeaders=" + strings.Join(signedHeaders, ";") +
		", Signature=" + signature
	request.Header.Set("Authorization", authorization)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request err: %w", err)
	}

	all, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return all, nil
}
