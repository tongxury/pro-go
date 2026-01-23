package biz

import (
	"context"
	userpb "store/api/user"
	typepb "store/api/user/types"
	"store/app/user/internal/data"
	"store/app/user/internal/data/enums"
	"store/app/user/internal/data/repo/ent"
	"store/app/user/internal/data/repo/ent/usercreditcost"
	"store/app/user/internal/data/repo/ent/usercreditincrement"
	"store/pkg/sdk/conv"
	"time"
)

type UserCreditBiz struct {
	data *data.Data
}

func NewUserCreditBiz(data *data.Data) *UserCreditBiz {
	return &UserCreditBiz{data: data}
}

func (t *UserCreditBiz) CostCredit(ctx context.Context, params *userpb.CheckCreditParams) (*userpb.CheckCreditResult, error) {

	summary, err := t.GetCreditSummary(ctx, params.UserId)
	if err != nil {
		return nil, err
	}

	if summary == nil || summary.Remaining < params.Cost {
		return &userpb.CheckCreditResult{
			Ok: false,
		}, nil
	}

	cost, err := t.data.Repos.EntClient.UserCreditCost.Create().
		SetUserID(params.UserId).
		SetAmount(params.Cost).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &userpb.CheckCreditResult{
		Ok:        true,
		CostId:    conv.Str(cost.ID),
		Remaining: summary.Remaining - params.Cost,
		Amount:    params.Cost,
	}, err
}

func (t *UserCreditBiz) Incr(ctx context.Context, userId string, amount int64) error {
	err := t.data.Repos.EntClient.UserCreditIncrement.Create().
		SetUserID(userId).
		SetAmount(amount).
		SetStatus(enums.UserCreditStatusActive).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *UserCreditBiz) ListOngoingPlanCredits(ctx context.Context, userId string) ([]*ent.UserCreditIncrement, error) {

	incrs, err := t.data.Repos.EntClient.UserCreditIncrement.Query().
		Where(usercreditincrement.UserID(userId)).
		Where(usercreditincrement.PlanIDNotNil()).
		Where(usercreditincrement.ExpireAtGTE(time.Now())).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return incrs, nil
}

func (t *UserCreditBiz) GetCreditSummary(ctx context.Context, userId string) (*typepb.CreditSummary, error) {

	var total []struct {
		UserId string `json:"user_id"`
		Amount int64  `json:"amount"`
	}

	err := t.data.Repos.EntClient.UserCreditIncrement.Query().
		Where(usercreditincrement.UserID(userId)).
		Where(usercreditincrement.Or(
			usercreditincrement.ExpireAtGTE(time.Now()),
			usercreditincrement.ExpireAtIsNil(),
		)).
		Where(usercreditincrement.Status(enums.UserCreditStatusActive)).
		GroupBy(usercreditincrement.FieldUserID).
		Aggregate(ent.As(ent.Sum(usercreditincrement.FieldAmount), "amount")).
		Scan(ctx, &total)

	if err != nil {
		return nil, err
	}

	if len(total) == 0 {
		return &typepb.CreditSummary{}, nil
	}

	// cost
	var cost int64
	var costs []struct {
		UserId string `json:"user_id"`
		Amount int64  `json:"amount"`
	}

	err = t.data.Repos.EntClient.UserCreditCost.Query().
		Where(usercreditcost.UserID(userId)).
		GroupBy(usercreditincrement.FieldUserID).
		Aggregate(ent.As(ent.Sum(usercreditincrement.FieldAmount), "amount")).
		Scan(ctx, &costs)

	if len(costs) > 0 {
		cost = costs[0].Amount
	}

	remaining := total[0].Amount - cost
	if remaining < 0 {
		remaining = 0
	}

	return &typepb.CreditSummary{
		Total:     total[0].Amount,
		Remaining: remaining,
	}, nil

}
