package tgbothelper

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"store/pkg/sdk/conv"
	"strings"
	"time"
)

func TgBotAuth(token string, isTest bool) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {

			request, ok := http.RequestFromServerContext(ctx)
			if !ok {
				return nil, errors.InternalServer("parse request err", "")
			}
			// We expect passing init data in the Authorization header in the following format:
			// <auth-type> <auth-data>
			// <auth-type> must be "tma", and <auth-data> is Telegram Mini Apps init data.
			authParts := strings.Split(request.Header.Get("Authorization"), " ")
			if len(authParts) != 2 {
				return nil, errors.Unauthorized("invalid authorization", "")
			}

			authType := authParts[0]
			if authType != "tma" {
				return nil, errors.Unauthorized("invalid authorization", "")
			}

			authData := authParts[1]

			if !isTest {
				//// Validate init data. We consider init data sign valid for 1 hour from their
				//// creation moment.
				if err := initdata.Validate(authData, token, 3000*time.Hour); err != nil {
					return nil, errors.Unauthorized("invalid authorization", "")
				}

			}

			// Parse init data. We will surely need it in the future.
			initData, err := initdata.Parse(authData)
			if err != nil {
				return nil, errors.Unauthorized("invalid authorization", "")
			}

			ctx = context.WithValue(ctx, "authData", &initData)

			//context.Request = context.Request.WithContext(
			//	withInitData(context.Request.Context(), initData),
			//)

			return handler(ctx, req)
		}
	}
}

type TgAuthData struct {
	*initdata.InitData
}

func (t *TgAuthData) UserID() string {
	return conv.String(t.User.ID)
}
func (t *TgAuthData) Nickname() string {
	return fmt.Sprintf("%s %s", t.User.FirstName, t.User.LastName)
}

func GetTgBotAuthData(ctx context.Context) *TgAuthData {

	return &TgAuthData{
		InitData: ctx.Value("authData").(*initdata.InitData),
	}
}
