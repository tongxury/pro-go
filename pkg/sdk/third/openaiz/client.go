package openaiz

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Client struct {
	c    openai.Client
	conf Config
}

type Config struct {
	AppKey  string
	BaseUrl string
	Proxy   string
}

func NewClient(conf Config) *Client {

	options := []option.RequestOption{
		option.WithAPIKey(conf.AppKey),
	}

	if conf.BaseUrl != "" {
		options = append(options, option.WithBaseURL(conf.BaseUrl))
	}

	// 添加 HTTP Client 配置（无论是否使用代理）
	var httpClient *http.Client

	if conf.Proxy != "" {
		proxyURL, err := url.Parse(conf.Proxy)
		if err != nil {
			panic(fmt.Sprintf("invalid proxy URL: %v", err))
		}

		// 完善的 Transport 配置
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			// 添加 TLS 配置以确保 HTTPS 连接正常
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // 生产环境建议设为 false
			},
		}

		httpClient = &http.Client{
			Transport: transport,
			Timeout:   3000 * time.Second, // 添加超时设置
		}
	} else {
		// 不使用代理时也创建一个标准的 HTTP Client
		httpClient = &http.Client{
			Timeout: 3000 * time.Second,
		}
	}

	// 始终添加 HTTP Client 选项
	options = append(options, option.WithHTTPClient(httpClient))

	client := openai.NewClient(
		options...,
	)

	//// 测试连接
	//resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
	//	Messages: []openai.ChatCompletionMessageParamUnion{
	//		openai.UserMessage("Say this is a test"),
	//	},
	//	Model: openai.ChatModelGPT4o,
	//})
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println(resp)

	return &Client{client, conf}
}

func (t Client) C() openai.Client {
	return t.c
}

func (t Client) Videos() *openai.VideoService {
	return &t.c.Videos
}

func (t Client) ChatCompletions() *openai.ChatCompletionService {
	return &t.c.Chat.Completions
}
