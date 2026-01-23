package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func main() {
	// 设置代理URL
	proxyURL, err := url.Parse("http://localhost:6060")
	if err != nil {
		panic(err)
	}

	// 创建带代理的HTTP客户端
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 30 * time.Second,
	}

	// 测试HTTP请求
	fmt.Println("Testing HTTP request through proxy...")
	resp, err := client.Get("http://httpbin.org/ip")
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("HTTP Response: %s\n", string(body))
	}

	// 测试HTTPS请求
	fmt.Println("\nTesting HTTPS request through proxy...")
	resp, err = client.Get("https://httpbin.org/ip")
	if err != nil {
		fmt.Printf("HTTPS request failed: %v\n", err)
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("HTTPS Response: %s\n", string(body))
	}

	// 测试代理统计信息
	fmt.Println("\nTesting proxy stats...")
	resp, err = client.Get("http://localhost:6060/stats")
	if err != nil {
		fmt.Printf("Stats request failed: %v\n", err)
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Proxy Stats: %s\n", string(body))
	}
}
