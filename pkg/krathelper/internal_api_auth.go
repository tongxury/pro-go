package krathelper

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	"strings"
)

func InternalAPIAuth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {

			request, ok := http.RequestFromServerContext(ctx)
			if !ok {
				return nil, errors.InternalServer("parse request err", "")
			}

			if strings.Contains(request.RequestURI, "/internal-api/") {
				apiKey := request.Header.Get("X-Api-Key")
				if apiKey != "urzZTmQZMvcxJHE44EOmjVwgTgr8Hqv4" {
					return nil, errors.Unauthorized("invalid api key", "")
				}
			}

			return handler(ctx, req)
		}
	}
}
