package repo

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"store/app/user/internal/data/repo/ent"
	"store/pkg/enums"
	"store/pkg/rediz"
	"time"
)

type memberSubscribeCache struct {
	redis *rediz.RedisClient
}

type CachedMemberSubscribe struct {
	ID            int64  `redis:"id"`
	UserID        int64  `redis:"user_id"`
	CreatedAt     int64  `redis:"created_at"`
	Level         string `redis:"level"`
	Status        string `redis:"status"`
	PaymentCycle  string `redis:"payment_cycle"`
	PromotionCode string `redis:"promotion_code"`
	ExpireAt      int64  `redis:"expire_at"`
	Uk            string `redis:"uk"`
}

func AsCached(sub *ent.MemberSubscribe) CachedMemberSubscribe {
	return CachedMemberSubscribe{
		ID:            sub.ID,
		UserID:        sub.UserID,
		CreatedAt:     sub.CreatedAt.Unix(),
		Level:         sub.Level.String(),
		Status:        sub.Status.String(),
		PaymentCycle:  sub.PaymentCycle.String(),
		PromotionCode: sub.PromotionCode,
		ExpireAt:      sub.ExpireAt.Unix(),
		Uk:            sub.Uk,
	}
}

func (t *CachedMemberSubscribe) AsEnt() *ent.MemberSubscribe {
	return &ent.MemberSubscribe{
		ID:            t.ID,
		UserID:        t.UserID,
		CreatedAt:     time.Unix(t.CreatedAt, 0),
		Level:         enums.MemberLevel(t.Level),
		Status:        enums.MemberSubscribeStatus(t.Status),
		PaymentCycle:  enums.PaymentCycle(t.PaymentCycle),
		PromotionCode: t.PromotionCode,
		ExpireAt:      time.Unix(t.ExpireAt, 0),
		Uk:            t.Uk,
	}
}

func (t *CachedMemberSubscribe) IsNilValue() bool {
	return t.ID == -1
}

var nilValue = CachedMemberSubscribe{ID: -1}

func (t *memberSubscribeCache) cacheKey(userID string) string {
	return fmt.Sprintf("member.subscribe.cachev2:%s", userID)
}
func (t *memberSubscribeCache) Find(ctx context.Context, userID string) (*CachedMemberSubscribe, bool, error) {

	var sub CachedMemberSubscribe
	err := t.redis.HGetAll(ctx, t.cacheKey(userID)).Scan(&sub)
	if err != nil {
		log.Errorw("HGetAll err", err, "userID", userID, "sub", sub)
		return nil, false, err
	}

	if sub.IsNilValue() {
		return nil, true, nil
	}

	if sub.ID == 0 {
		return nil, false, nil
	}

	return &sub, false, nil
}

func (t *memberSubscribeCache) SetNil(ctx context.Context, userID string, expire time.Duration) error {

	err := t.redis.HSet(ctx, t.cacheKey(userID), map[string]interface{}{
		"id": -1,
	}).Err()
	if err != nil {
		log.Errorw("HSet err", err, "userID", userID)
		return err
	}
	t.redis.Expire(ctx, t.cacheKey(userID), expire)

	return nil
}

func (t *memberSubscribeCache) Set(ctx context.Context, userID string, sub CachedMemberSubscribe, expire time.Duration) error {

	err := t.redis.HSet(ctx, t.cacheKey(userID), sub).Err()
	if err != nil {
		log.Errorw("HSet err", err, "userID", userID, "sub", sub)
		return err
	}

	t.redis.Expire(ctx, t.cacheKey(userID), expire)

	return nil
}

func (t *memberSubscribeCache) Delete(ctx context.Context, userID string) error {
	err := t.redis.Del(ctx, t.cacheKey(userID)).Err()
	if err != nil {
		log.Errorw("Del err", err, "userID", userID)
		return err
	}
	return nil
}
