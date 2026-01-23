package repo

import (
	"store/app/payment/internal/data/repo/ent"
	"store/pkg/rediz"
)

type PaymentRepo struct {
	db    *ent.Client
	redis *rediz.RedisClient
}

func NewPaymentRepo(db *ent.Client, redis *rediz.RedisClient) *PaymentRepo {
	return &PaymentRepo{db: db, redis: redis}
}
