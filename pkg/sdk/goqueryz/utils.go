package goqueryz

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Meta struct {
	Name    string
	Content string
}

func FindMetaContents(selection *goquery.Document, name string) []string {
	// 存储所有匹配的 meta content 值
	var contents []string

	// 定义递归查找函数
	var findAllMeta func(*html.Node)
	findAllMeta = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var nameMatches bool
			var content string

			// 检查该节点的所有属性
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == name {
					nameMatches = true
				}

				if attr.Key == "content" {
					content = attr.Val
				}
			}

			// 如果找到匹配的name和content，添加到结果中
			if nameMatches && content != "" {
				contents = append(contents, content)
			}
		}

		// 递归检查所有子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findAllMeta(c)
		}
	}

	// 处理文档中的所有顶级节点
	for _, node := range selection.Nodes {
		findAllMeta(node)
	}

	return contents
}

func FindMetaContent(selection *goquery.Document, name string) string {

	var findMeta func(*html.Node) string
	findMeta = func(n *html.Node) string {

		if n.Type == html.ElementNode && n.Data == "meta" {

			var ok bool

			for _, attr := range n.Attr {

				if attr.Key == "name" && attr.Val == name {
					ok = true
				}

				if ok && attr.Key == "content" {
					return attr.Val
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			v := findMeta(c)
			if v != "" {
				return v
			}
		}

		return ""
	}

	return findMeta(selection.Nodes[0])

}
