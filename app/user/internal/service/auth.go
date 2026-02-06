package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	responsepb "store/api/public/response"
	userpb "store/api/user"
	typepb "store/api/user/types"
	"store/app/user/internal/biz"
	"store/app/user/internal/data"
	"store/app/user/internal/data/repo"
	pkgConf "store/pkg/confcenter"
	"store/pkg/events"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/mathz"
	"store/pkg/sdk/mail"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	userpb.UnimplementedAuthServer
	authBiz *biz.AuthBiz
	credit  *biz.UserCreditBiz
	data    *data.Data
}

func NewAuthService(authBiz *biz.AuthBiz, credit *biz.UserCreditBiz, data *data.Data) *AuthService {
	return &AuthService{
		authBiz: authBiz,
		credit:  credit,
		data:    data,
	}
}

func (t *AuthService) GetAuthTokenV2(ctx context.Context, params *userpb.GetAuthTokenV2Request) (*userpb.AuthToken, error) {
	deviceId := krathelper.GetDeviceId(ctx)
	channel := krathelper.GetChannel(ctx)
	inviterCode := krathelper.GetQuery(ctx, "ic")

	token, userId, isRegister, err := t.authBiz.GetAuthTokenV2(ctx, params.GetCredential(), channel, deviceId, params.From, inviterCode)
	log.Errorw("GetAuthTokenV2", params.From)
	if err != nil {
		return nil, err
	}

	//cookieDomain, err := s.authBiz.GetCookieDomain()
	//if err != nil {
	//	return nil, err
	//}

	return &userpb.AuthToken{
		Token:  token,
		IsNew:  isRegister,
		UserId: userId,
	}, err
}

func (t *AuthService) GetEmailAuthCode(ctx context.Context, params *userpb.GetEmailAuthCodeParams) (*emptypb.Empty, error) {

	redisKey := fmt.Sprintf("authCodeEmail:%s", params.Email)

	code := mathz.RandNumber(100000, 999999)

	log.Debugw("验证码", "", "email", params.Email, "code", code)

	err := t.data.MailClient.Send(mail.Params{
		Subject: "【Veogo】登录验证",
		To:      params.Email,
		Content: fmt.Sprintf("<div>【Veogo】您的验证码为:%d, 5分钟之内有效。</div>", code),
	})
	if err != nil {
		return nil, err
	}

	err = t.data.Repos.RedisClient.Set(ctx, redisKey, code, 5*time.Minute).Err()
	if err != nil {
		log.Errorw("RedisClient set error", err, "params", params)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (t *AuthService) GetEmailAuthToken(ctx context.Context, params *userpb.GetEmailAuthTokenParams) (*userpb.AuthToken, error) {
	deviceId := krathelper.GetDeviceId(ctx)
	channel := krathelper.GetChannel(ctx)

	redisKey := fmt.Sprintf("authCodeEmail:%s", params.Email)

	code, err := t.data.Repos.RedisClient.Get(ctx, redisKey).Result()
	if err != nil {

		if errors.Is(err, redis.Nil) {
			return nil, errors.BadRequest("invalidCode", "")
		}

		log.Errorw("RedisClient get error", err, "params", params)
		return nil, err
	}

	if code != params.Code {
		return nil, errors.Unauthorized("invalidCode", "验证码错误")
	}

	id, isNew, err := t.data.Repos.User.InsertNXByEmail(ctx, repo.InsertNXByEmailParams{
		Email:   params.Email,
		Channel: channel,
	})

	//id, err := t.data.Repos.EntClient.User.Create().
	//	SetKey(params.Phone).
	//	SetStatus(enums.UserStatusActive).
	//	SetPhone(params.Phone).
	//	OnConflictColumns(user.FieldKey).
	//	DoNothing().ExecX(ctx)

	if err != nil {
		log.Errorw("UserCreate err", err, "params", params)
		return nil, err
	}

	log.Debugw("UserCreate id", id, "params", params)

	token, err := krathelper.GenerateToken(id)

	log.Debugw("UserCreate id", id, "params", params, "token", token)

	if isNew {
		t.sendXhsTrack(ctx, params.XhsClickId)
	}

	msg := kafka.Message{
		Value: conv.S2B(events.AuthEvent{
			UserID:     conv.Str(id),
			LoginBy:    "email",
			TS:         time.Now().Unix(),
			DeviceID:   deviceId,
			IsRegister: isNew,
		}),
	}

	err = t.data.KafkaClient.W().Write(ctx, events.Topic_AuthLogin, msg)

	if err != nil {
		log.Errorw("KafkaClient.W().Write err", err, "msg", msg)
	}

	return &userpb.AuthToken{
		Token:  token,
		IsNew:  isNew,
		UserId: conv.Str(id),
	}, nil
}

func (t *AuthService) sendXhsTrack(ctx context.Context, clickId string) {
	helper.Go(ctx, func(ctx context.Context) {
		err := t.data.XhsClient.Send(ctx, "7688267", clickId, "102")
		if err != nil {
			log.Errorw("SendXhsClickId err", err, "clickId", clickId)
			return
		}
		log.Debugw("SendXhsClickId id", clickId, "clickId", clickId)
	})
}

func (t *AuthService) GetPhoneAuthCode(ctx context.Context, params *userpb.GetPhoneAuthCodeParams) (*emptypb.Empty, error) {

	redisKey := fmt.Sprintf("authCode:%s", params.Phone)

	code := mathz.RandNumber(100000, 999999)

	log.Debugw("验证码", "", "phone", params.Phone, "code", code)

	_, err := t.data.AlismsClient.Send(ctx, []string{params.Phone}, conv.Str(code))
	if err != nil {
		return nil, err
	}

	err = t.data.Repos.RedisClient.Set(ctx, redisKey, code, 5*time.Minute).Err()
	if err != nil {
		log.Errorw("RedisClient set error", err, "params", params)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (t *AuthService) GetPhoneAuthToken(ctx context.Context, params *userpb.GetPhoneAuthTokenParams) (*userpb.AuthToken, error) {
	deviceId := krathelper.GetDeviceId(ctx)

	channel := krathelper.GetChannel(ctx)

	redisKey := fmt.Sprintf("authCode:%s", params.Phone)

	code, err := t.data.Repos.RedisClient.Get(ctx, redisKey).Result()
	if err != nil {
		log.Errorw("RedisClient get error", err, "params", params)
		return nil, err
	}

	if code != params.Code {
		return nil, errors.Unauthorized("", "验证码错误")
	}

	id, isNew, err := t.data.Repos.User.InsertNX(ctx, params.Phone, channel)

	//id, err := t.data.Repos.EntClient.User.Create().
	//	SetKey(params.Phone).
	//	SetStatus(enums.UserStatusActive).
	//	SetPhone(params.Phone).
	//	OnConflictColumns(user.FieldKey).
	//	DoNothing().ExecX(ctx)

	if err != nil {
		log.Errorw("UserCreate err", err, "params", params)
		return nil, err
	}

	token, err := krathelper.GenerateToken(id)

	if isNew {
		t.sendXhsTrack(ctx, params.XhsClickId)
	}

	msg := kafka.Message{
		Value: conv.S2B(events.AuthEvent{
			UserID:     conv.Str(id),
			LoginBy:    "phone",
			TS:         time.Now().Unix(),
			DeviceID:   deviceId,
			IsRegister: isNew,
		}),
	}

	err = t.data.KafkaClient.W().Write(ctx, events.Topic_AuthLogin, msg)

	if err != nil {
		log.Errorw("KafkaClient.W().Write err", err, "msg", msg)
	}

	return &userpb.AuthToken{
		Token:  token,
		IsNew:  isNew,
		UserId: conv.Str(id),
	}, nil
}

func (t *AuthService) GetWxAuthToken(ctx context.Context, params *userpb.GetWxAuthTokenParams) (*userpb.AuthToken, error) {

	deviceId := krathelper.GetDeviceId(ctx)

	number, err := t.data.WXDev.GetUserPhoneNumber(ctx, params.Code)
	if err != nil {
		return nil, err
	}

	id, isNew, err := t.data.Repos.User.InsertNX(ctx, number, "")

	if err != nil {
		log.Errorw("UserCreate err", err, "params", params)
		return nil, err
	}

	tokenStr, err := krathelper.GenerateToken(id)
	if err != nil {
		return nil, err
	}

	msg := kafka.Message{
		Value: conv.S2B(events.AuthEvent{
			UserID:     conv.Str(id),
			LoginBy:    "apple",
			TS:         time.Now().Unix(),
			DeviceID:   deviceId,
			IsRegister: isNew,
		}),
	}

	err = t.data.KafkaClient.W().Write(ctx, events.Topic_AuthLogin, msg)

	if err != nil {
		log.Errorw("KafkaClient.W().Write err", err, "msg", msg)
	}

	return &userpb.AuthToken{
		Token:  tokenStr,
		IsNew:  isNew,
		UserId: conv.Str(id),
	}, nil

}

func (t *AuthService) GetAppleAuthToken(ctx context.Context, params *userpb.GetAppleAuthTokenParams) (*userpb.AuthToken, error) {

	deviceId := krathelper.GetDeviceId(ctx)
	channel := krathelper.GetChannel(ctx)

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

	id, isNew, err := t.data.Repos.User.InsertNXByEmail(ctx, repo.InsertNXByEmailParams{
		Email:    claims["email"].(string),
		Channel:  channel,
		Nickname: params.Nickname,
		Platform: "apple",
	})

	if err != nil {
		log.Errorw("UserCreate err", err, "params", params)
		return nil, err
	}

	tokenStr, err := krathelper.GenerateToken(id)
	if err != nil {
		return nil, err
	}

	msg := kafka.Message{
		Value: conv.S2B(events.AuthEvent{
			UserID:     conv.Str(id),
			LoginBy:    "apple",
			TS:         time.Now().Unix(),
			DeviceID:   deviceId,
			IsRegister: isNew,
		}),
	}

	err = t.data.KafkaClient.W().Write(ctx, events.Topic_AuthLogin, msg)

	if err != nil {
		log.Errorw("KafkaClient.W().Write err", err, "msg", msg)
	}

	return &userpb.AuthToken{
		Token:  tokenStr,
		IsNew:  isNew,
		UserId: conv.Str(id),
	}, nil
}

func (t *AuthService) GetAuthToken(ctx context.Context, req *userpb.GetAuthTokenRequest) (*responsepb.RedirectResponse, error) {

	channel, platform, err := t.authBiz.ParseState(req.GetState())
	if err != nil {
		return nil, err
	}

	//deviceID :=  kratosutil.GetHeader(ctx, "Device-Id")

	token, _, err := t.authBiz.GetAuthToken(ctx, req.GetCode(), channel, channel)
	if err != nil {
		return nil, err
	}

	redirectTo := t.authBiz.GetRedirectTo(platform)

	cookieDomain, err := t.authBiz.GetCookieDomain()
	if err != nil {
		return nil, err
	}

	//go func() {
	//	helper.DeferFunc()
	//
	//	err := s.eventLogBiz.AddEventLog(context.Background(), &adminpb.AddEventLogParams{
	//		EventId: "google_login_success_v2",
	//		UserId:  conv.Int64(userId),
	//	})
	//	if err != nil {
	//		log.Errorw("AddEventLog err", err)
	//	}
	//}()

	return &responsepb.RedirectResponse{
		Cookies: []*responsepb.Cookie{
			{Name: pkgConf.AuthCookieName, Value: token, Domain: cookieDomain},
		},
		Url: redirectTo,
	}, err
}

func (t *AuthService) GetAuthCode(ctx context.Context, req *userpb.GetAuthCodeRequest) (*responsepb.RedirectResponse, error) {

	platform := req.GetPlatform()
	if platform == "" {
		platform = "web"
	}

	authPlatform := req.GetAuthPlatform()
	if authPlatform == "" {
		authPlatform = "google"
	}

	url, err := t.authBiz.GetAuthCodeUrl(ctx, req.GetPlatform(), req.GetAuthPlatform(), req.GetChannel())
	if err != nil {
		return nil, err
	}

	return &responsepb.RedirectResponse{Url: url}, nil

}

func (t *AuthService) GetAuthUser(ctx context.Context, req *userpb.GetAuthUserRequest) (*typepb.User, error) {

	userId := krathelper.RequireUserId(ctx)

	dbUser, err := t.authBiz.GetAuthUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &typepb.User{
		//Id:     conv.String(dbUser.ID),
		Username: dbUser.Nickname,
		//UserAvatar: conv.String(dbUser.GoogleAuth["picture"]),
		Email: dbUser.Email,
	}, nil

}

func (t *AuthService) AddRegisterUser(ctx context.Context, params *userpb.AddRegisterUserParams) (*emptypb.Empty, error) {
	deviceID := krathelper.GetDeviceId(ctx)

	message, err := t.authBiz.AddAuthUser(ctx, params.Email, params.Pwd, params.Nickname, deviceID)
	if err != nil {
		log.Errorw("AddAuthUser err", err, "params", params)
		return nil, err
	}

	if message != "" {
		return nil, errors.BadRequest("", message)
	}

	return &emptypb.Empty{}, nil
}

func (t *AuthService) GetVerifyCode(ctx context.Context, params *userpb.GetVerifyCodeParams) (*emptypb.Empty, error) {

	message, err := t.authBiz.GetVerifyCode(ctx, params.Email, params.Type)
	if err != nil {
		log.Errorw("AddAuthUser err", err, "params", params)
		return nil, err
	}

	if message != "" {
		return nil, errors.BadRequest("", message)
	}

	return &emptypb.Empty{}, nil
}

func (t *AuthService) SubmitVerifyCode(ctx context.Context, params *userpb.SubmitVerifyCodeParams) (*emptypb.Empty, error) {

	failReason, err := t.authBiz.VerifyRegister(ctx, params.Email, params.Code)
	if err != nil {
		log.Errorw("VerifyRegister err", err, "params", params)
		return nil, err
	}

	if failReason != "" {
		return nil, errors.BadRequest("", failReason)
	}

	return &emptypb.Empty{}, nil
}

func (t *AuthService) GetAuthTokenByEmail(ctx context.Context, params *userpb.GetAuthTokenByEmailParams) (*userpb.GetAuthTokenByEmailResult, error) {

	deviceID := krathelper.GetDeviceId(ctx)

	result, reason, err := t.authBiz.GetAuthTokenByEmail(ctx, params.Email, params.Pwd, deviceID)
	if err != nil {
		log.Errorw("GetAuthTokenByEmail err", err, "params", params)
		return nil, err
	}

	if result == nil {
		return nil, errors.Unauthorized("", reason)
	}

	return result, nil
}

func (t *AuthService) ResetAuthPassword(ctx context.Context, params *userpb.ResetAuthPasswordParams) (*emptypb.Empty, error) {

	err := t.authBiz.ResetPassword(ctx, params.Email, params.Pwd)
	if err != nil {
		log.Errorw("ResetPassword err", err, "params", params)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
