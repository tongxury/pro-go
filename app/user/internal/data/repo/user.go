package repo

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/patrickmn/go-cache"
	"store/app/user/internal/data/enums"
	"store/app/user/internal/data/repo/ent"
	"store/app/user/internal/data/repo/ent/user"
	"store/pkg/middlewares/entz"
	"store/pkg/rediz"
	"store/pkg/sdk/conv"
	"time"
)

type UserRepo struct {
	db    *ent.Client
	redis *rediz.RedisClient
	lc    *cache.Cache
}

func NewUserRepo(db *ent.Client, redis *rediz.RedisClient, lcClient *cache.Cache) *UserRepo {
	return &UserRepo{
		db:    db,
		redis: redis,
		lc:    lcClient,
	}
}

type InsertNXByEmailParams struct {
	Email    string
	Channel  string
	Nickname string
	Avatar   string
	Platform string
}

func (t *UserRepo) InsertNXByEmail(ctx context.Context, params InsertNXByEmailParams) (int64, bool, error) {

	olds, err := t.db.User.Query().Where(user.Key(params.Email)).All(ctx)
	if err != nil {
		return 0, false, err
	}

	if len(olds) > 0 {
		return olds[0].ID, false, nil
	}

	//key := helper.OrString(params.Key, params.Email)

	id, err := t.db.User.Create().
		SetKey(params.Email).
		SetStatus(enums.UserStatusActive).
		SetEmail(params.Email).
		SetNickname(params.Nickname).
		SetChannel(params.Channel).
		SetAvatar(params.Avatar).
		SetPlatform(params.Platform).
		OnConflictColumns(user.FieldKey).
		DoNothing().
		ID(ctx)
	if err != nil {
		return 0, false, err
	}

	return id, true, nil

}

func (t *UserRepo) InsertNX(ctx context.Context, phone, channel string) (int64, bool, error) {

	olds, err := t.db.User.Query().Where(user.Key(phone)).All(ctx)
	if err != nil {
		return 0, false, err
	}

	if len(olds) > 0 {
		return olds[0].ID, false, nil
	}

	id, err := t.db.User.Create().
		SetKey(phone).
		SetStatus(enums.UserStatusActive).
		SetPhone(phone).
		OnConflictColumns(user.FieldKey).
		DoNothing().
		ID(ctx)
	if err != nil {
		return 0, false, err
	}

	return id, true, nil

}

func (t *UserRepo) Upsert(ctx context.Context, u *ent.User) (*ent.User, error) {

	userId, err := t.db.User.Create().
		SetNickname(u.Nickname).
		SetEmail(u.Email).
		SetKey(u.Email).
		//OnConflict(sql.ConflictColumns("google_auth_user_id")).
		OnConflict(sql.ConflictColumns(user.FieldKey)).
		Update(func(upsert *ent.UserUpsert) {
		}).
		ID(ctx)

	if err != nil {
		return nil, err
	}

	// todo 目前没有好的办法获取 执行后的最终数据状态
	newUser, err := t.db.User.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (t *UserRepo) FindById(ctx context.Context, userId string) (*User, error) {

	dbUsers, err := t.db.User.Query().
		Where(user.ID(conv.Int64(userId))).
		All(ctx)

	if err != nil {
		return nil, err
	}
	if len(dbUsers) == 0 {
		return nil, nil
	}

	return &User{
		User: *dbUsers[0],
	}, nil
}

func (t *UserRepo) GetById(ctx context.Context, userId string, useCache bool) (*User, error) {

	key := fmt.Sprintf("user.cache:%s", userId)

	if useCache {
		if x, found := t.lc.Get(key); found {
			return x.(*User), nil
		}
	}

	x, err := t.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if x == nil {
		return nil, errors.NotFound("no user found by id : "+userId, "")
	}

	t.lc.Set(key, x, time.Minute)

	return x, nil
}

type ListParams struct {
	Keyword      string
	UserID       string
	Emails       []string
	EmailOrPhone string
	Page, Size   int64
}

func (t *UserRepo) List(ctx context.Context, params ListParams) ([]*User, int64, error) {

	q := t.db.User.Query()

	if params.UserID != "" {
		q = q.Where(user.ID(conv.Int64(params.UserID)))
	}

	if params.Keyword != "" {
		q = q.Where(user.Or(
			user.ID(conv.Int64(params.Keyword)),
			user.EmailHasPrefix(params.Keyword),
			user.NicknameHasPrefix(params.Keyword),
			user.PhoneHasPrefix(params.Keyword),
		))
	}

	if params.EmailOrPhone != "" {
		q = q.Where(user.Or(
			user.Email(params.EmailOrPhone),
			user.Phone(params.EmailOrPhone),
		))
	}

	if len(params.Emails) > 0 {
		q = q.Where(user.EmailIn(params.Emails...))
	}

	count, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset, limit, _, use := entz.VerifyPageSize(params.Page, params.Size)
	if use {
		q = q.Limit(limit).Offset(offset)
	}

	users, err := q.Order(ent.Desc(user.FieldID)).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return t.asUsers(users), int64(count), nil
}

func (t *UserRepo) asUsers(users []*ent.User) []*User {
	var rsp []*User
	for _, user := range users {
		rsp = append(rsp, &User{*user})
	}
	return rsp
}
