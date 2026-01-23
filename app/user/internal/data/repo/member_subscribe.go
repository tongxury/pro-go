package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"store/app/user/internal/data/repo/ent"
	"store/app/user/internal/data/repo/ent/membersubscribe"
	"store/pkg/enums"
	"store/pkg/middlewares/entz"
	"store/pkg/rediz"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/types"
	"time"
)

type MemberSubscribeRepo struct {
	db *ent.Client
	//redis *clients.RedisClusterClient
	//lc    *cache.Cache
	cache *memberSubscribeCache
}

func NewMemberSubscribeRepo(db *ent.Client, redis *rediz.RedisClient) *MemberSubscribeRepo {
	return &MemberSubscribeRepo{db: db, cache: &memberSubscribeCache{redis: redis}}
}

type ListMemberSubscribesParams struct {
	UserIDs    []int64
	Expired    bool
	Page, Size int64
}

func (t *MemberSubscribeRepo) ListMemberSubscribes(ctx context.Context, params ListMemberSubscribesParams) (MemberSubscribes, int64, error) {

	q := t.db.MemberSubscribe.Query()

	if len(params.UserIDs) > 0 {
		q = q.Where(membersubscribe.UserIDIn(params.UserIDs...))
	}

	if params.Expired {
		q = q.Where(membersubscribe.ExpireAtLT(time.Now()))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset, limit, _, use := entz.VerifyPageSize(params.Page, params.Size)
	if use {
		q = q.Limit(offset).Offset(limit)
	}

	dbMemberSubs, err := q.Order(ent.Desc(membersubscribe.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	return AsMemberSubscribes(dbMemberSubs), int64(total), nil
}

func (t *MemberSubscribeRepo) FindLatestMemberSubscribe(ctx context.Context, userID string, useCache bool) (*ent.MemberSubscribe, error) {

	if useCache {
		sub, nilInDB, err := t.cache.Find(ctx, userID)

		if err == nil {
			if nilInDB {
				return nil, nil
			}

			if sub != nil {
				return sub.AsEnt(), nil
			}
		}
	}

	dbMemberSubs, err := t.db.MemberSubscribe.Query().
		Where(membersubscribe.UserID(conv.Int64(userID))).
		Order(ent.Desc(membersubscribe.FieldID)).
		Limit(1).
		All(ctx)

	if err != nil {
		return nil, err
	}

	if len(dbMemberSubs) == 0 {
		_ = t.cache.SetNil(ctx, userID, 5*time.Minute)
		return nil, nil
	}
	_ = t.cache.Set(ctx, userID, AsCached(dbMemberSubs[0]), 5*time.Minute)

	return dbMemberSubs[0], nil
}

func (t *MemberSubscribeRepo) ClearSubscribe(ctx context.Context, userID string) error {

	_ = t.cache.Delete(ctx, userID)

	_, err := t.db.MemberSubscribe.Delete().
		Where(membersubscribe.UserID(conv.Int64(userID))).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *MemberSubscribeRepo) UpsertSubscribe(ctx context.Context, userID, level, cycle, outSubId, version, promotionCode string, expireAt int64) (int64, int, error) {

	tx, err := t.db.Tx(ctx)
	if err != nil {
		return 0, 0, err
	}

	uk := outSubId
	if version != "" {
		uk = version + "_" + outSubId
	}

	oldSubs, err := tx.MemberSubscribe.Query().Where(membersubscribe.Uk(uk)).All(ctx)
	if err != nil {
		return 0, 0, err
	}

	var isRenew bool
	var id int64
	if len(oldSubs) > 0 {

		if oldSubs[0].ExpireAt.Unix() == expireAt {
			return 0, 0, nil
		}

		err = tx.MemberSubscribe.Update().
			SetExpireAt(time.Unix(expireAt, 0)).
			SetStatus(enums.MemberSubscribeStatus_Subscribing).
			Where(membersubscribe.ID(oldSubs[0].ID)).Exec(ctx)
		isRenew = true
		id = oldSubs[0].ID
	} else {
		id, err = tx.MemberSubscribe.Create().
			SetUserID(conv.Int64(userID)).
			SetLevel(enums.MemberLevel(level)).
			SetPaymentCycle(enums.PaymentCycle(cycle)).
			SetStatus(enums.MemberSubscribeStatus_Subscribing).
			SetExpireAt(time.Unix(expireAt, 0)).
			SetPromotionCode(promotionCode).
			SetExtra(types.MemberSubscribeExtra{OutSubId: outSubId}).
			SetUk(uk).
			OnConflictColumns(membersubscribe.FieldUk).
			Update(func(upsert *ent.MemberSubscribeUpsert) {
				upsert.SetExpireAt(time.Unix(expireAt, 0))
			}).
			ID(ctx)
	}

	if err != nil {
		return 0, 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, err
	}

	_ = t.cache.Delete(ctx, userID)
	log.Debugw("Delete cache", "", "userID", userID)

	return id, helper.Select(isRenew, 2, 1), nil
}
