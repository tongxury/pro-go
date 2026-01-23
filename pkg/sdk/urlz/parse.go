package urlz

import (
	"net/url"
	"strings"
)

type URLInfo struct {
	Protocol     string            `json:"protocol"`
	Host         string            `json:"host"`
	Hostname     string            `json:"hostname"`
	Port         string            `json:"port"`
	Path         string            `json:"path"`
	PathSegments []string          `json:"pathSegments"`
	QueryString  string            `json:"queryString"`
	QueryParams  map[string]string `json:"queryParams"`
	Fragment     string            `json:"fragment"`
	FullURL      string            `json:"fullUrl"`
}

func ParseURL(urlString string) (*URLInfo, error) {
	// 解析URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	// 解析查询参数
	queryParams := make(map[string]string)
	for k, v := range parsedURL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0] // 取第一个值
		}
	}

	// 解析路径段
	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	// 过滤空字符串
	var segments []string
	for _, segment := range pathSegments {
		if segment != "" {
			segments = append(segments, segment)
		}
	}

	return &URLInfo{
		Protocol:     parsedURL.Scheme,
		Host:         parsedURL.Host,
		Hostname:     parsedURL.Hostname(),
		Port:         parsedURL.Port(),
		Path:         parsedURL.Path,
		PathSegments: segments,
		QueryString:  parsedURL.RawQuery,
		QueryParams:  queryParams,
		Fragment:     parsedURL.Fragment,
		FullURL:      urlString,
	}, nil
}
