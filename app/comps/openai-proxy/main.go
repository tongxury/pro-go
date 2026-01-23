package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort = "6000"
	defaultURL  = "https://api.openai.com/v1"
)

// OpenAIProxy 纯转发代理，不做任何数据处理
type OpenAIProxy struct {
	backendURL string
	client     *http.Client // 使用默认的超时配置
	logger     *log.Logger
}

// NewOpenAIProxy 创建代理实例
func NewOpenAIProxy() *OpenAIProxy {
	backendURL := os.Getenv("OPENAI_BACKEND_URL")
	if backendURL == "" {
		backendURL = defaultURL
	}

	return &OpenAIProxy{
		backendURL: backendURL,
		// 使用默认的 DefaultClient，它使用 DefaultTransport
		// DefaultTransport 没有全局超时限制，适合透传
		client: http.DefaultClient,
		logger: log.New(os.Stdout, "[OpenAI-Proxy] ", log.LstdFlags),
	}
}

// ServeHTTP 纯转发处理
func (p *OpenAIProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 记录日志
	p.logger.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	// 构建完整后端URL
	backendURL := p.backendURL + r.URL.Path
	if r.URL.RawQuery != "" {
		backendURL += "?" + r.URL.RawQuery
	}

	// 创建新请求，保持上下文
	req, err := http.NewRequestWithContext(r.Context(), r.Method, backendURL, r.Body)
	if err != nil {
		p.logger.Printf("Error creating request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 拷贝所有请求头（除了 Host，让它自动填充）
	for key, values := range r.Header {
		if key != "Host" {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	// 转发请求
	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Printf("Error forwarding request: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 拷贝所有响应头
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 设置状态码
	w.WriteHeader(resp.StatusCode)

	// 直接流式拷贝响应体
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		// 写入失败通常是因为客户端断开连接，不需要额外处理
		return
	}

	p.logger.Printf("%s %s -> %d", r.Method, r.URL.Path, resp.StatusCode)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	proxy := NewOpenAIProxy()

	// 使用默认配置，不设置超时
	server := &http.Server{
		Addr:    ":" + port,
		Handler: proxy,
	}

	proxy.logger.Printf("Starting OpenAI Proxy")
	proxy.logger.Printf("Listening on :%s", port)
	proxy.logger.Printf("Backend: %s", proxy.backendURL)
	proxy.logger.Printf("Set OPENAI_BACKEND_URL to customize backend")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		proxy.logger.Fatalf("Server failed: %v", err)
	}
}
