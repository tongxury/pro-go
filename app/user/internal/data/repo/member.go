package repo

import (
	"context"
	errors2 "errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"store/app/user/internal/data/repo/ent"
	"store/pkg/rediz"
	"store/pkg/sdk/conv"
	"store/pkg/types"
	"time"
)

type MemberRepo struct {
	db    *ent.Client
	redis *rediz.RedisClient
}

func NewMemberRepo(db *ent.Client, redis *rediz.RedisClient) *MemberRepo {
	return &MemberRepo{db: db, redis: redis}
}

func (t *MemberRepo) IncrUsage(ctx context.Context, userId, model, function string, modelLimit *types.Limit, functionLimit *types.Limit) (bool, bool, error) {

	now := time.Now()
	mm := now.Format("2006-01")
	dd := now.Format("2006-01-02")

	key := fmt.Sprintf("user.member.used.v9:%s", userId)

	fields := []string{
		model,
		fmt.Sprintf("%s:%s", model, mm),
		fmt.Sprintf("%s:%s", model, dd),
	}

	if function != "" {
		fields = append(fields,
			function,
			fmt.Sprintf("%s:%s", function, mm),
			fmt.Sprintf("%s:%s", function, dd))
	}

	var useBonus bool

	_, err := t.redis.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {

		usages := make(Usage, len(fields))

		_ = pipeliner.HMGet(ctx, key, fields...)

		oldCmds, err := pipeliner.Exec(ctx)
		if err != nil {
			return err
		}

		oldValues := oldCmds[0].(*redis.SliceCmd).Val()

		for i, field := range fields {
			usages[field] = conv.Int64(oldValues[i]) + 1
		}

		// model限制
		modelAvailable := true
		if modelLimit != nil {
			modelAvailable = modelLimit.Check(
				usages.FindDayUsage(model),
				usages.FindMonthUsage(model),
				usages.FindTotalUsage(model),
			)
			if !modelAvailable {
				modelAvailable = modelLimit.Bonus > 0
				if modelAvailable {
					useBonus = true
				}

			}
		}

		// 功能限制
		functionAvailable := true
		if functionLimit != nil {
			functionAvailable = functionLimit.Check(
				usages.FindDayUsage(function),
				usages.FindMonthUsage(function),
				usages.FindTotalUsage(function),
			)
		}

		if !modelAvailable || !functionAvailable {
			return Limited{}
		}

		for _, field := range fields {
			pipeliner.HIncrBy(ctx, key, field, 1)
		}

		_, err = pipeliner.Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	if err == nil {
		return true, useBonus, nil
	}

	if errors.Is(err, Limited{}) {
		return false, false, nil
	}

	return false, false, err
}

func (t *MemberRepo) AddMemberUsed(ctx context.Context, userId, model, function string) error {

	now := time.Now()
	mm := now.Format("2006-01")
	dd := now.Format("2006-01-02")

	key := fmt.Sprintf("user.member.used.v8:%s", userId)

	fields := []string{
		model,
		fmt.Sprintf("%s:%s", model, mm),
		fmt.Sprintf("%s:%s", model, dd),
	}

	if function != "" {
		fields = append(fields,
			function,
			fmt.Sprintf("%s:%s", function, mm),
			fmt.Sprintf("%s:%s", function, dd))
	}

	log.Debugf("AddMemberUsed %v", fields)

	_, err := t.redis.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {

		for _, field := range fields {
			if err := pipeliner.HIncrBy(ctx, key, field, 1).Err(); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *MemberRepo) GetMemberUsed(ctx context.Context, userId string) (map[string]int64, error) {

	//now := time.Now()

	key := fmt.Sprintf("user.member.used.v9:%s", userId)

	result, err := t.redis.HGetAll(ctx, key).Result()
	if err != nil && !errors2.Is(err, redis.Nil) {
		return nil, err
	}

	rsp := make(map[string]int64, len(result))
	for k, v := range result {
		rsp[k] = conv.Int64(v)
	}

	return rsp, nil
}
