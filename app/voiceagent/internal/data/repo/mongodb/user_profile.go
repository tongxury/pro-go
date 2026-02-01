package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserProfileCollection struct {
	*mgz.Core[*voiceagent.UserProfile]
}

func NewUserProfileCollection(db *mongo.Database) *UserProfileCollection {
	return &UserProfileCollection{
		Core: mgz.NewCore[*voiceagent.UserProfile](db, "va_user_profiles"),
	}
}
