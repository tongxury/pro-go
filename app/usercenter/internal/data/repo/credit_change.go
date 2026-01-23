package repo

import (
	creditpb "store/api/credit"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type CreditChangeCollection struct {
	*mgz.Core[*creditpb.CreditChange]
}

func NewCreditChangeCollection(db *mongo.Database) *CreditChangeCollection {
	return &CreditChangeCollection{
		Core: mgz.NewCore[*creditpb.CreditChange](db, "credit_changes"),
	}
}
