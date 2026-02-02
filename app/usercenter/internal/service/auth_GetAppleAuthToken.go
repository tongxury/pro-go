package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	ucpb "store/api/usercenter"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"
)

func (t AuthService) GetAppleAuthToken(ctx context.Context, params *ucpb.GetAppleAuthTokenRequest) (*ucpb.Token, error) {

	set, err := jwk.Fetch(context.Background(), "https://appleid.apple.com/auth/keys")

	if err != nil {
		return nil, err
	}

	// 解析 token
	token, err := jwt.Parse(params.IdentityToken, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		key, found := set.LookupKeyID(kid)
		if !found {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		var publicKey rsa.PublicKey
		if err := key.Raw(&publicKey); err != nil {
			return nil, fmt.Errorf("failed to parse public key: %v", err)
		}

		return &publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// 验证 token 是否过期
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("token expired")
	}

	//// 验证 audience (你的 Apple Service ID)
	////host.exp.Exponent
	//if !claims.VerifyAudience("host.exp.Exponent", true) {
	//	return nil, fmt.Errorf("invalid audience")
	//}

	// 验证 issuer
	if !claims.VerifyIssuer("https://appleid.apple.com", true) {
		return nil, fmt.Errorf("invalid issuer")
	}

	userId, isNew, err := t.data.Mongo.User.InsertNX(ctx,
		&ucpb.User{
			Key:       params.Email,
			Email:     params.Email,
			CreatedAt: time.Now().Unix(),
		},
		mgz.Filter().EQ("key", params.Email).
			B(),
	)

	if err != nil {
		log.Errorw("UserCreate err", err, "params", params)
		return nil, err
	}

	tokenStr, err := krathelper.GenerateTokenV2(userId)
	if err != nil {
		return nil, err
	}

	return &ucpb.Token{
		Token:  tokenStr,
		IsNew:  isNew,
		UserId: userId,
	}, nil
}
