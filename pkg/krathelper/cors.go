package krathelper

import (
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"net/http"
	"strconv"
	"strings"
)

var defaultCORSHeaders = []string{"Referer", "Content-Type", HEADER_NAME_AUTHORIZATION,
	HEADER_NAME_PLATFORM, HEADER_NAME_DEVICE_ID, HEADER_NAME_CHANNEL, HEADER_NAME_LOCATION, HEADER_NAME_LOCALE, HEADER_NAME_PROMOTION_CODE}

// CORS 配置结构体
type corsConfig struct {
	allowedHeaders   []string
	allowedMethods   []string
	allowedOrigins   []string
	allowCredentials bool
	maxAge           int
}

// CORS 中间件
func corsMiddleware(cfg corsConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 设置允许的源
			origin := r.Header.Get("Origin")
			if origin == "" {
				origin = "*"
			}

			// 设置 CORS 头
			w.Header().Set("Access-Control-Allow-Origin", origin)

			if cfg.allowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.allowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.allowedHeaders, ","))
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.maxAge))
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// DefaultCorsConfig 默认的 CORS 配置
var DefaultCorsConfig = khttp.Filter(corsMiddleware(corsConfig{
	allowedHeaders:   defaultCORSHeaders,
	allowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	allowedOrigins:   []string{"*"},
	maxAge:           2592000,
	allowCredentials: true,
}))

// CorsConfig 自定义 CORS 配置
func CorsConfig(headers ...string) khttp.ServerOption {
	return khttp.Filter(corsMiddleware(corsConfig{
		allowedHeaders:   append(defaultCORSHeaders, headers...),
		allowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		allowedOrigins:   []string{"*"},
		maxAge:           2592000,
		allowCredentials: true,
	}))
}

// 可选：添加一个更灵活的配置构建器
type CORSOption func(*corsConfig)

func WithHeaders(headers []string) CORSOption {
	return func(c *corsConfig) {
		c.allowedHeaders = headers
	}
}

func WithMethods(methods []string) CORSOption {
	return func(c *corsConfig) {
		c.allowedMethods = methods
	}
}

func WithMaxAge(seconds int) CORSOption {
	return func(c *corsConfig) {
		c.maxAge = seconds
	}
}

// NewCorsConfig 使用选项模式创建 CORS 配置
func NewCorsConfig(opts ...CORSOption) khttp.ServerOption {
	cfg := corsConfig{
		allowedHeaders:   defaultCORSHeaders,
		allowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		allowedOrigins:   []string{"*"},
		maxAge:           600,
		allowCredentials: true,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return khttp.Filter(corsMiddleware(cfg))
}
