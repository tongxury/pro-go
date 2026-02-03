package krathelper

import (
	"context"
	"errors"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

var SecretSignKey = "1447cdd12451aceb97cc0ea8e5c1cc906f94dbb05c2397abcee8b64345a6b56e"

func GenerateToken(userId int64) (string, error) {

	claims := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, jwtv4.MapClaims{
		"user_id": strconv.FormatInt(userId, 10),
		"ts":      time.Now().Unix(),
	})
	signedString, err := claims.SignedString([]byte(SecretSignKey))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func GenerateTokenV2(userId string) (string, error) {

	claims := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, jwtv4.MapClaims{
		"user_id": userId,
		"ts":      time.Now().Unix(),
	})
	signedString, err := claims.SignedString([]byte(SecretSignKey))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func ParseUserID(token string) (string, error) {
	jwtToken, err := jwtv4.Parse(token, func(token *jwtv4.Token) (interface{}, error) {
		return []byte(SecretSignKey), nil
	})
	if err != nil {
		return "", err
	}

	cc := jwtToken.Claims
	ccc := cc.(jwtv4.MapClaims)

	userID, found := ccc["user_id"]

	if found {
		return userID.(string), err
	}

	//mapClaims := claims.(jwtv4.MapClaims)

	return "", nil
}

func ParseClaims(token string, key string) (*jwtv4.Claims, error) {

	jwtToken, err := jwtv4.Parse(token, func(token *jwtv4.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	return &jwtToken.Claims, nil
}

func RequireUserId(ctx context.Context) string {
	claims, _ := jwt.FromContext(ctx)
	return claims.(jwtv4.MapClaims)["user_id"].(string)
}

func RequireUserID(ctx context.Context) int64 {
	return conv.Int64(RequireUserId(ctx))
}

func GetUserId(ctx context.Context) (string, error) {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return "", errors.New("no claims found in context")
	}

	userInfo, ok := claims.(jwtv4.MapClaims)

	log.Debugw("userInfo", userInfo)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userID, found := userInfo["user_id"]
	if !found {
		return "", errors.New("no user_id found")
	}

	return userID.(string), nil
}

func FindUserId(ctx context.Context) string {
	claims, ok := jwt.FromContext(ctx)
	if ok {
		return claims.(jwtv4.MapClaims)["user_id"].(string)
	}
	return ""
}

func Auth(whiteList []string) middleware.Middleware {
	return selector.Server(
		jwt.Server(
			func(token *jwtv4.Token) (interface{}, error) {
				return []byte(SecretSignKey), nil
			},
			//jwt.WithSigningMethod(jwtv4.SigningMethodHS256),
			//jwt.WithClaims(func() jwtv4.Claims {
			//	resp := &jwtv4.MapClaims{}
			//	fmt.Println(resp)
			//	return resp
			//}),
		),
	).Match(func(ctx context.Context, operation string) bool {
		return !helper.InSlice(operation, whiteList)
	}).Build()
}

func NormalizeAuthorization(key string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			request, ok := http.RequestFromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			token := request.Header.Get("Authorization")

			if token == "" {
				token = request.URL.Query().Get("token")
			}

			if token == "" {
				token = request.URL.Query().Get("authorization")
			}

			if token != "" {

				claims, err := ParseClaims(token, key)
				if err != nil {
					return nil, err
				}

				if !strings.HasPrefix(token, "Bearer") {
					token = "Bearer " + token
				}

				ctx = jwt.NewContext(ctx, *claims)
				request.Header.Set("Authorization", token)
			}

			return handler(ctx, req)
		}
	}
}

func FromCookie(name string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			request, ok := http.RequestFromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			cookie, err := request.Cookie(name)
			if err != nil {
				return handler(ctx, req)
			}

			if cookie != nil {
				request.Header.Set("Authorization", "Bearer "+cookie.Value)
			}

			return handler(ctx, req)
		}
	}
}

func FindToken() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			request, ok := http.RequestFromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			tokenName := "authToken"

			token := request.Header.Get("Authorization")
			if token == "" {
				token = request.Header.Get(tokenName)
			}

			if token == "" {
				cookie, _ := request.Cookie(tokenName)
				if cookie != nil {
					token = cookie.Value
				}
			}

			if token == "" {
				token = request.URL.Query().Get(tokenName)
			}

			request.Header.Set("Authorization", "Bearer "+token)

			return handler(ctx, req)
		}
	}
}
