package repo

import (
	ucpb "store/api/usercenter"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserCollection struct {
	*mgz.Core[*ucpb.User]
}

func NewUserCollection(db *mongo.Database) *UserCollection {
	return &UserCollection{
		Core: mgz.NewCore[*ucpb.User](db, "users"),
	}
}
