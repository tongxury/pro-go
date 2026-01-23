package repo

import (
	"context"
	"store/pkg/clients"
	"time"
)

type StripeRepo struct {
	db *clients.ClickHouseClient
}

func NewStripeRepo(db *clients.ClickHouseClient) *StripeRepo {
	return &StripeRepo{db: db}
}

func (t *StripeRepo) Insert(ctx context.Context, e *StripeBill) error {

	err := t.db.Exec(ctx, `
		insert into stripe_bills (
		    id, user_id, created_at, member_level, member_cycle, account_country, billing_reason, currency, serial, status, 
		    subscription_id, subtotal,
		    customer_country, customer_email, customer_name,
		    promotion_code, coupon_id, coupon_name, coupon_percent_off, coupon_times_redeemed, v) 
				values (?,?,?,?,?,?,?,?,?,?,
				        ?,?,
				        ?,?,?,
				        ?,?,?,?,?,?
				)`,
		e.ID, e.UserID, e.CreatedAt.Format(time.DateTime), e.MemberLevel, e.MemberCycle, e.AccountCountry, e.BillingReason, e.Currency, e.Serial, e.Status,
		e.SubscriptionID, e.Subtotal,
		e.CustomerCountry, e.CustomerEmail, e.CustomerName,
		e.PromotionCode, e.CouponID, e.CouponName, e.CouponPercentOff, e.CouponTimesRedeemed, e.Version,
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *StripeRepo) FindLastByCreatedAt(ctx context.Context) (*time.Time, error) {

	var result struct {
		CreatedAt *time.Time `ch:"created_at"`
	}
	err := t.db.QueryRow(ctx, "select max(created_at) as created_at from stripe_bills").ScanStruct(&result)
	if err != nil {
		return nil, err
	}

	return result.CreatedAt, nil
}

func (t *StripeRepo) FindRetryingBills(ctx context.Context) ([]StripeBill, error) {

	var results []StripeBill
	err := t.db.Select(ctx, &results, "select id, v from stripe_bills where status not in ('paid', 'void')")
	if err != nil {
		return nil, err
	}

	return results, nil
}
