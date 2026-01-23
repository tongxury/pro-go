package krathelper

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"strings"
)

const (
	HEADER_NAME_LOCALE         = "Locale"
	HEADER_NAME_LOCATION       = "Location"
	HEADER_NAME_CHANNEL        = "Channel"
	HEADER_NAME_DEVICE_ID      = "Device-Id"
	HEADER_NAME_PLATFORM       = "Platform"
	HEADER_NAME_AUTHORIZATION  = "Authorization"
	HEADER_NAME_PROMOTION_CODE = "Promotion-Code"
	//HEADER_NAME_CLIENT_VERSION = "Client-Version"
)

func GetChannel(ctx context.Context) string {
	channel := GetHeader(ctx, "Channel")
	if channel != "" {
		return channel
	}

	return GetQuery(ctx, "cn")
}

func GetDeviceId(ctx context.Context) string {
	channel := GetHeader(ctx, "Device-Id")
	if channel != "" {
		return channel
	}

	return GetQuery(ctx, "dv")
}

func GetHeaderLocale(ctx context.Context) string {
	return GetHeader(ctx, HEADER_NAME_LOCALE)
}

func FindHeader(ctx context.Context, name string) string {
	return GetHeader(ctx, name)
}

func GetHeader(ctx context.Context, name string) string {

	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return ""
	}

	return tr.RequestHeader().Get(name)
}

func GetQuery(ctx context.Context, name string) string {

	request, ok := http.RequestFromServerContext(ctx)
	if !ok {
		return ""
	}

	q := request.URL.RawQuery
	if q == "" {
		return ""
	}

	for _, x := range strings.Split(q, "&") {
		kv := strings.Split(x, "=")
		if len(kv) != 2 {
			continue
		}
		if kv[0] == name {
			return kv[1]
		}
	}

	return ""
}

func FindCookie(ctx context.Context, name string) string {

	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return ""
	}

	cookieString := tr.RequestHeader().Get("Cookie")

	for _, x := range strings.Split(cookieString, ";") {

		parts := strings.Split(x, "=")

		if len(parts) == 2 && parts[0] == name {
			return parts[1]
		}

	}

	return ""
}

func NormalizeContentType() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			request, ok := http.RequestFromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			//log.Debugw("NormalizeContentType", request.Header.Get("Content-Type"))

			c := request.Header.Get("Content-Type")
			if c == "" {
				request.Header.Set("Content-Type", "application/json")
			}

			//if strings.HasPrefix(c, "multipart/form-data") {
			//	// 去掉 boundary=xxx 相关参数
			//	request.Header.Set("Content-Type", "multipart/form-data")
			//}
			//
			//log.Debugw("NormalizeContentType", request.Header.Get("Content-Type"))

			return handler(ctx, req)
		}
	}
}
