package biz

import (
	"context"
	errors2 "errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"io"
	mgmtpb "store/api/mgmt"
	userpb "store/api/user"
	"store/app/user/internal/data"
	"store/app/user/internal/data/enums"
	"store/app/user/internal/data/repo"
	"store/app/user/internal/data/repo/ent"
	"store/app/user/internal/data/repo/ent/user"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/rand"

	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"strings"
	"time"
)

type AuthBiz struct {
	data *data.Data
}

func NewAuthBiz(data *data.Data) *AuthBiz {
	return &AuthBiz{data: data}
}

func (t *AuthBiz) getUserInfoByCredential(ctx context.Context, credential string) (map[string]interface{}, error) {

	var userInfo map[string]interface{}

	resp, err := resty.New().R().SetContext(ctx).Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + credential)
	if err != nil {
		return nil, err
	}

	err = conv.B2S(resp.Body(), &userInfo)
	if err != nil {
		return nil, err
	}

	if len(userInfo) == 0 {
		return nil, fmt.Errorf("no user found by credential: %s", credential)
	}

	return userInfo, nil
}

func (t *AuthBiz) GetAuthTokenV2(ctx context.Context, credential, channel, deviceID, from, inviterCode string) (string, string, bool, error) {
	userInfo, err := t.getUserInfoByCredential(ctx, credential)
	if err != nil {
		return "", "", false, err
	}

	log.Debugw("get user info", "", "userInfo", userInfo)

	id, b, err := t.data.Repos.User.InsertNXByEmail(ctx, repo.InsertNXByEmailParams{
		Email:    conv.String(userInfo["email"]),
		Channel:  channel,
		Nickname: conv.String(userInfo["nickname"]),
		Avatar:   conv.String(userInfo["picture"]),
	})
	if err != nil {
		return "", "", false, err
	}

	token, _ := krathelper.GenerateToken(id)

	return token, conv.String(id), b, nil
}

func (t *AuthBiz) GetAuthUser(ctx context.Context, userId string) (*ent.User, error) {

	u, err := t.data.Repos.EntClient.User.Query().
		Where(user.ID(conv.Int64(userId))).
		First(ctx)

	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, fmt.Errorf("no user found by id: %s")
	}

	return u, nil
}

func (t *AuthBiz) GetAuthToken(ctx context.Context, code, channel, deviceID string) (string, string, error) {

	userInfo, err := t.getAuthToken(ctx, code)
	if err != nil {
		return "", "", err
	}

	p := &ent.User{
		Nickname: conv.String(userInfo["name"]),
		Email:    conv.String(userInfo["email"]),
	}

	newUser, err := t.data.Repos.User.Upsert(ctx, p)
	if err != nil {
		return "", "", err
	}

	//userId, err := t.data.Repos.EntClient.User.Create().
	//	SetGoogleAuth(userInfo).
	//	SetSerial(serial).
	//	SetUID(deviceID).
	//	SetChannel(channel).
	//	SetInviteCode(serial).
	//	SetBindCode(channel).
	//	SetNickname(conv.String(userInfo["name"])).
	//	SetEmail(conv.String(userInfo["email"])).
	//	SetAvatar(conv.String(userInfo["picture"])).
	//	//OnConflict(sql.ConflictColumns("google_auth_user_id")).
	//	OnConflict(sql.ConflictColumns(userauth.FieldEmail)).
	//	UpdateFields(func(upsert *ent.UserUpsert) {
	//		upsert.UpdateGoogleAuth()
	//		upsert.UpdateAvatar()
	//		upsert.Set(userauth.FieldLastLoginTime, time.Now())
	//	}).
	//	ID(ctx)
	//
	//if err != nil {
	//	return "", "", err
	//}
	token, err := krathelper.GenerateToken(newUser.ID)
	if err != nil {
		return "", "", err
	}

	//err = t.data.Repos.RedisClient.XAdd(ctx, &redis.XAddArgs{
	//	Stream: confcenter.Topic_AuthLogin,
	//	Values: events.NewAuthEvent(
	//		conv.String(userId),
	//		confcenter.AuthLoginBy_Google,
	//		deviceID,
	//	),
	//}).Err()
	//if err != nil {
	//	log.Errorw("XAdd err", err, "userID", userId, "LoginBy", confcenter.AuthLoginBy_Google)
	//}
	//
	//
	//

	//err = t.data.KafkaClient.W().Write(ctx, events.Topic_AuthLogin, kafka.Message{
	//	Value: conv.S2B(events.AuthEvent{
	//		UserID:   newUser.ID,
	//		LoginBy:  confcenter.AuthLoginBy_Google,
	//		TS:       time.Now().Unix(),
	//		DeviceID: deviceID,
	//	}),
	//})
	//
	//if err != nil {
	//	log.Errorw("KafkaClient Write err", err, "userId", newUser.ID)
	//}

	return token, conv.String(newUser.ID), nil

}

func (t *AuthBiz) getAuthToken(ctx context.Context, code string) (map[string]interface{}, error) {

	token, err := t.data.BizConfig.Oauth2.Google.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := t.data.BizConfig.Oauth2.Google.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}

	userInfoBytes, _ := io.ReadAll(resp.Body)

	err = conv.B2S(userInfoBytes, &userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (t *AuthBiz) GetAuthCodeUrl(ctx context.Context, platform, authPlatform, channel string) (string, error) {
	url := t.data.BizConfig.Oauth2.Google.AuthCodeURL(t.generateState(channel, platform),
		oauth2.SetAuthURLParam("prompt", "select_account"),
	)
	return url, nil
}

func (t *AuthBiz) generateState(channel, platform string) string {
	return fmt.Sprintf("%s.%s.%s", shortuuid.New(), channel, platform)
}

func (t *AuthBiz) ParseState(state string) (string, string, error) {
	parts := strings.Split(state, ".")

	if len(parts) != 3 {
		return "", "", errors.BadRequest("invalid state: "+state, "")
	}

	return parts[1], parts[2], nil
}

func (t *AuthBiz) GetRedirectTo(platform string) string {

	if platform == "" {
		platform = "web"
	}

	redirectTo := t.data.BizConfig.Oauth2.Redirect[platform]

	return redirectTo
}

func (t *AuthBiz) GetCookieDomain() (string, error) {
	domain := t.data.BizConfig.Oauth2.CookieDomain

	return domain, nil
}

func (t *AuthBiz) ResetPassword(ctx context.Context, email, pwd string) error {
	// todo
	err := t.data.Repos.EntClient.User.Update().
		Where(user.Email(email)).
		SetPassword(helper.MD5([]byte(pwd))).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *AuthBiz) GetAuthTokenByEmail(ctx context.Context, email, pwd, deviceID string) (*userpb.GetAuthTokenByEmailResult, string, error) {

	users, err := t.data.Repos.EntClient.User.Query().
		Where(user.Email(email)).Where(user.Status(enums.UserStatusActive)).
		All(ctx)
	if err != nil {
		return nil, "", err
	}

	if len(users) == 0 {
		return nil, "The user does not exists", nil
	}

	if users[0].Password != helper.MD5([]byte(pwd)) {
		return nil, "your password is wrong, please retry", nil
	}

	err = t.data.Repos.EntClient.User.Update().
		Where(user.Email(email)).
		Exec(ctx)
	if err != nil {
		return nil, "", err
	}

	token, err := krathelper.GenerateToken(users[0].ID)
	if err != nil {
		return nil, "", err
	}

	//err = t.data.Repos.RedisClient.XAdd(ctx, &redis.XAddArgs{
	//	Stream: confcenter.Topic_AuthLogin,
	//	Values: events.NewAuthEvent(
	//		users[0].ID,
	//		confcenter.AuthLoginBy_Email,
	//		deviceID,
	//	),
	//}).Err()
	//if err != nil {
	//	log.Errorw("XAdd err", err, "userID", users[0].ID, "LoginBy", confcenter.AuthLoginBy_Email)
	//}

	//msg := kafka.Message{
	//	Value: conv.S2B(events.AuthEvent{
	//		UserID:   users[0].ID,
	//		LoginBy:  confcenter.AuthLoginBy_Email,
	//		TS:       time.Now().Unix(),
	//		DeviceID: deviceID,
	//	}),
	//}
	//
	//err = t.data.KafkaClient.W().Write(ctx, events.Topic_AuthLogin, msg)
	//
	//if err != nil {
	//	log.Errorw("KafkaClient.W().Write err", err, "msg", msg)
	//}
	return &userpb.GetAuthTokenByEmailResult{
		Token:    token,
		Redirect: t.GetRedirectTo(""),
	}, "", nil
}

func (t *AuthBiz) GetVerifyCode(ctx context.Context, email string, typ int64) (string, error) {

	//verifyState := 0
	//if typ == 2 { // 忘记密码
	//	verifyState = 1
	//}

	users, err := t.data.Repos.EntClient.User.Query().
		Where(user.Email(email)).Where(user.Status(enums.UserStatusPending)).
		All(ctx)
	if err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "The email is not found, please register first", nil
	}

	code := rand.Range(100000, 999999)

	_, err = t.data.Repos.RedisClient.Set(ctx, "reg_code:"+email, code, 5*time.Minute).Result()
	if err != nil {
		return "", err
	}

	_, err = t.data.GrpcClients.MgmtClient.SendEmail(ctx, &mgmtpb.SendEmailParams{
		Recipients: []*mgmtpb.SendEmailParams_Recipient{
			{
				Email:      email,
				Subject:    "Verify your email ID",
				TemplateId: "2",
				Args: &mgmtpb.SendEmailParams_Args{
					Username:   users[0].Nickname,
					VerifyCode: conv.String(code),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	return "", nil

}

func (t *AuthBiz) VerifyRegister(ctx context.Context, email, code string) (string, error) {

	savedCode, err := t.data.Repos.RedisClient.Get(ctx, "reg_code:"+email).Result()
	if err != nil && !errors2.Is(err, redis.Nil) {
		return "", err
	}

	if savedCode == "" {
		return "verify code has been expired", nil
	}

	if code != savedCode {
		return "verify code is wrong", nil
	}

	err = t.data.Repos.EntClient.User.Update().
		SetStatus(enums.UserStatusActive).
		Where(user.Email(email)).Exec(ctx)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (t *AuthBiz) AddAuthUser(ctx context.Context, email, pwd, nickname, deviceID string) (string, error) {

	exists, err := t.data.Repos.EntClient.User.Query().
		Where(user.Email(email)).Where(user.Status(enums.UserStatusPending)).
		Exist(ctx)
	if err != nil {
		return "", err
	}

	if exists {
		return "The user already exists, please log in", nil
	}

	err = t.data.Repos.EntClient.User.Create().
		SetEmail(email).
		SetPassword(helper.MD5([]byte(pwd))).
		SetNickname(nickname).
		OnConflictColumns(user.FieldEmail).
		Update(func(upsert *ent.UserUpsert) {
			upsert.UpdateNickname()
			upsert.UpdatePassword()
		}).
		Exec(ctx)
	if err != nil {
		return "", err
	}

	return "", nil
}
