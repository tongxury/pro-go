package repo

import (
	creditpb "store/api/credit"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type CreditRechargeCollection struct {
	*mgz.Core[*creditpb.CreditRecharge]
}

func NewCreditRechargeCollection(db *mongo.Database) *CreditRechargeCollection {
	return &CreditRechargeCollection{
		Core: mgz.NewCore[*creditpb.CreditRecharge](db, "credit_recharges"),
	}
}
