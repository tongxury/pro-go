package geminiai

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"store/pkg/sdk/helper/mathz"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/genai"
	"google.golang.org/grpc"
)

type GenaiFactory struct {
	clients []*Client
}

type GeminiProxyRoundTripper struct {
	APIKey string
	//ProxyURL  string
	Transport http.RoundTripper
}

func (t *GeminiProxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// 克隆请求
	newReq := req.Clone(req.Context())

	// 添加 API 密钥到 header（而不是查询参数）
	newReq.Header.Set("x-goog-api-key", t.APIKey)

	return t.Transport.RoundTrip(newReq)
}

func newClient(config Config) (*Client, error) {

	options := []option.ClientOption{
		option.WithAPIKey(config.Key),
		option.WithGRPCDialOption(grpc.WithIdleTimeout(600 * time.Second)),
		//option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
		//	Time:                600 * time.Second,
		//	Timeout:             600 * time.Second,
		//	PermitWithoutStream: true,
		//})),
	}

	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   600 * time.Second, // 增加连接超时
				KeepAlive: 600 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       600 * time.Second, // 增加空闲连接超时
			TLSHandshakeTimeout:   600 * time.Second, // 增加TLS握手超时
			ExpectContinueTimeout: 600 * time.Second, // 增加Expect超时
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			// 添加响应头超时
			ResponseHeaderTimeout: 600 * time.Second,
		}
		// 更新选项
		c := &http.Client{
			//Transport: &http.Transport{
			//	Proxy: http.ProxyURL(proxyURL),
			//},
			Transport: &GeminiProxyRoundTripper{
				APIKey: config.Key,
				//ProxyURL:  config.Proxy,
				Transport: transport,
			},
		}

		options = []option.ClientOption{
			option.WithAPIKey(config.Key),
			option.WithHTTPClient(c),
		}
	}

	c, err := genai.NewClient(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		c: c,
	}, nil
}

func NewGenaiFactory(config FactoryConfig) *GenaiFactory {

	var clients []*Client

	for i := range config.Configs {
		genaiClient, err := newClient(config.Configs[i])
		if err != nil {
			panic(err)
		}

		//if i == 0 {
		//	model := genaiClient.c.GenerativeModel("gemini-2.5-flash-lite-preview-06-17")
		//	resp, err := model.GenerateContent(context.Background(), genai.Text("Hello"))
		//	if err != nil {
		//		panic(err)
		//	}
		//	log.Debugw("genai resp", resp)
		//}

		clients = append(clients, genaiClient)
	}

	return &GenaiFactory{
		clients: clients,
	}
}

func (t GenaiFactory) Get() *Client {
	idx := mathz.RandNumber(0, len(t.clients)-1)
	return t.clients[idx]
}
