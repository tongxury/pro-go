package biz

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	paymentpb "store/api/payment"
	"store/app/payment/internal/data"
	"store/app/payment/internal/data/repo/ent"
	"store/app/payment/internal/data/repo/ent/creditrecharge"
	"store/pkg/sdk/conv"
	"strings"
	"time"
)

type CreditBiz struct {
	data *data.Data
}

func NewCreditBiz(data *data.Data) *CreditBiz {
	return &CreditBiz{data: data}
}

func (t *CreditBiz) Cost(ctx context.Context, userId, key string, cost int64) error {

	state, err := t.GetCreditState(ctx, userId)
	if err != nil {
		return err
	}

	if state.Remaining < cost {
		return errors.New("exceeded")
	}

	tx, err := t.data.Repos.EntClient.Tx(ctx)
	if err != nil {
		return err
	}

	err = tx.CreditRecharge.Update().
		AddCost(cost).
		Where(creditrecharge.ID(conv.Int64(state.RechargeId))).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = tx.CreditConsume.Create().
		SetKey(key).
		SetRechargeID(state.RechargeId).
		SetUserID(userId).
		SetAmount(cost).
		Exec(ctx)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (t *CreditBiz) GetCreditState(ctx context.Context, userId string) (*paymentpb.CreditState, error) {

	//crs, err := t.data.Repos.EntClient.CreditRecharge.Query().
	//	Where(creditrecharge.UserID(userId)).
	//	Where(creditrecharge.StatusEQ("completed")).
	//	Where(creditrecharge.ExpireAtGT(time.Now())).
	//	Order(ent.Desc(creditrecharge.FieldExpireAt)).
	//	All(ctx)

	crs, err := t.data.Repos.EntClient.CreditRecharge.Query().
		Where(creditrecharge.UserID(userId)).
		Where(creditrecharge.StatusEQ("completed")).
		Where(creditrecharge.ExpireAtGT(time.Now())).
		Where(func(s *sql.Selector) {
			s.Where(sql.GT(s.C(creditrecharge.FieldQuota), s.C(creditrecharge.FieldCost)))
		}).
		Order(ent.Desc(creditrecharge.FieldExpireAt)).
		All(ctx)

	if err != nil {
		return nil, err
	}

	log.Debugw("CreditRecharge crs", conv.S2J(crs))

	if len(crs) == 0 {
		return &paymentpb.CreditState{
			Show: strings.HasSuffix(userId, "1") ||
				strings.HasSuffix(userId, "2") ||
				strings.HasSuffix(userId, "3") ||
				strings.HasSuffix(userId, "4"),
		}, nil
	}

	total := crs[0].Quota
	cost := crs[0].Cost
	// 把所有的 credit余额都汇集到最后一条充值记录 简化逻辑
	if len(crs) > 1 {

		tx, err := t.data.Repos.EntClient.Tx(ctx)
		if err != nil {
			return nil, err
		}

		var extra int64
		for i := 1; i < len(crs); i++ {
			x := crs[i]

			err = tx.CreditRecharge.Update().
				Where(creditrecharge.ID(x.ID)).
				SetQuota(x.Cost).
				Exec(ctx)

			if err != nil {
				return nil, err
			}

			extra += x.Quota - x.Cost
		}

		err = tx.CreditRecharge.Update().
			Where(creditrecharge.ID(crs[0].ID)).
			AddQuota(extra).
			Exec(ctx)

		if err != nil {
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			return nil, err
		}

		total += extra
	}

	return &paymentpb.CreditState{
		UserId:     userId,
		Total:      total,
		Used:       cost,
		Remaining:  total - cost,
		RechargeId: conv.Str(crs[0].ID),
		ExpireAt:   crs[0].ExpireAt.Unix(),
	}, nil
}
