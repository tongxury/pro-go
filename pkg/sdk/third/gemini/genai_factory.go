package gemini

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"store/pkg/sdk/helper/mathz"
	"time"

	"cloud.google.com/go/auth/credentials"
	"cloud.google.com/go/auth/httptransport"
	"google.golang.org/genai"
)

type GenaiFactory struct {
	clients []*Client
}

type ProxyRoundTripper struct {
	APIKey string
	//ProxyURL  string
	Transport http.RoundTripper
}

func (t *ProxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// 克隆请求
	newReq := req.Clone(req.Context())

	// 添加 API 密钥到 header（而不是查询参数）
	newReq.Header.Set("x-goog-api-key", t.APIKey)

	return t.Transport.RoundTrip(newReq)
}

func newClient(config *Config) (*Client, error) {

	ctx := context.Background()

	var httpClient *http.Client

	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   600 * time.Second,
				KeepAlive: 600 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       600 * time.Second,
			TLSHandshakeTimeout:   600 * time.Second,
			ExpectContinueTimeout: 600 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			ResponseHeaderTimeout: 600 * time.Second,
		}

		// Use ProxyRoundTripper for API Key mode
		if config.Key != "" {
			httpClient = &http.Client{
				Transport: &ProxyRoundTripper{
					APIKey:    config.Key,
					Transport: transport,
				},
			}
		} else {
			// For Vertex AI with Credentials, we need to add auth middleware
			httpClient = &http.Client{
				Transport: transport,
			}
		}
	}

	clientConfig := &genai.ClientConfig{
		HTTPClient: httpClient,
		HTTPOptions: genai.HTTPOptions{
			BaseURL: config.BaseURL,
		},
	}

	if config.Project != "" {
		clientConfig.Backend = genai.BackendVertexAI
		clientConfig.Project = config.Project
		clientConfig.Location = config.Location

		if config.CredentialsJSON != "" {
			cred, err := credentials.DetectDefault(&credentials.DetectOptions{
				Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
				CredentialsJSON: []byte(config.CredentialsJSON),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create credentials from JSON: %w", err)
			}
			// Explicitly add auth middleware to the custom HTTP Client if it exists
			if httpClient != nil {
				err = httptransport.AddAuthorizationMiddleware(httpClient, cred)
				if err != nil {
					return nil, fmt.Errorf("failed to add auth middleware: %w", err)
				}
			}
		}
	} else {
		// Only set APIKey if not using Vertex AI (or if strictly using API Key mode)
		clientConfig.APIKey = config.Key
	}

	client, err := genai.NewClient(ctx, clientConfig)

	if err != nil {
		return nil, err
	}

	return &Client{
		c:     client,
		cache: config.Cache,
	}, nil
}

func NewGenaiFactory(config *FactoryConfig) *GenaiFactory {

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
