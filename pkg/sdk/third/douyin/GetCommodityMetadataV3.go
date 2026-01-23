package douyin

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/html"
	"strings"
	"time"
)

func (t *Client) GetCommodityMetadataV3(ctx context.Context, url string) (*CommodityMetadata, error) {
	// 创建resty客户端
	client := resty.New()

	// 设置客户端配置
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)

	// 设置请求头
	client.SetHeaders(map[string]string{
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept-Encoding":           "gzip, deflate, br",
		"Connection":                "keep-alive",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Cache-Control":             "max-age=0",
	})

	// 发送请求
	resp, err := client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 检查响应状态码
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode())
	}

	// 获取响应内容
	html := resp.String()
	fmt.Println(html)

	return nil, nil
}

// 辅助函数：根据选择器查找节点
func findNodeBySelector(n *html.Node, selector string) *html.Node {
	var result *html.Node
	var traverse func(*html.Node)

	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode {
			// 检查class属性
			for _, attr := range node.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, strings.TrimPrefix(selector, ".")) {
					result = node
					return
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(n)
	return result
}

// 辅助函数：根据选择器查找多个节点
func findNodesBySelector(n *html.Node, selector string) []*html.Node {
	var results []*html.Node
	var traverse func(*html.Node)

	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode {
			// 检查class属性
			for _, attr := range node.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, strings.TrimPrefix(selector, ".")) {
					results = append(results, node)
					break
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(n)
	return results
}

// 辅助函数：提取文本内容
func extractText(n *html.Node) string {
	var text strings.Builder
	var traverse func(*html.Node)

	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			text.WriteString(strings.TrimSpace(node.Data))
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(n)
	return strings.TrimSpace(text.String())
}
