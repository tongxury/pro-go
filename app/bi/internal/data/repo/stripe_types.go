package repo

import "time"

type StripeBill struct {
	ID             string    `ch:"id"`
	UserID         int64     `ch:"user_id"`
	MemberLevel    string    `ch:"member_level"`
	MemberCycle    string    `ch:"member_cycle"`
	AccountCountry string    `ch:"account_country"`
	BillingReason  string    `ch:"billing_reason"`
	CreatedAt      time.Time `ch:"created_at"`
	Currency       string    `ch:"currency"`
	Serial         string    `ch:"serial"`
	Status         string    `ch:"status"`
	Subtotal       int64     `ch:"subtotal"`

	SubscriptionID string `ch:"subscription_id"`

	CustomerCountry string `ch:"customer_country"`
	CustomerEmail   string `ch:"customer_email"`
	CustomerName    string `ch:"customer_name"`

	PromotionCode       string  `ch:"promotion_code"`
	CouponID            string  `ch:"coupon_id"`
	CouponName          string  `ch:"coupon_name"`
	CouponPercentOff    float64 `ch:"coupon_percent_off"`
	CouponTimesRedeemed int64   `ch:"coupon_times_redeemed"`
	Version             string  `ch:"v"`
}
