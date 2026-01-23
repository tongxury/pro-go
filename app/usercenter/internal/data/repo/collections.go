package repo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	Database       *mongo.Database
	User           *UserCollection
	CreditRecharge *CreditRechargeCollection
	CreditChange   *CreditChangeCollection
}

func NewCollections(database *mongo.Database) *Collections {
	return &Collections{
		Database:       database,
		User:           NewUserCollection(database),
		CreditRecharge: NewCreditRechargeCollection(database),
		CreditChange:   NewCreditChangeCollection(database),
	}
}
