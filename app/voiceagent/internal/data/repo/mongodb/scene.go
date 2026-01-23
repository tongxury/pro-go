package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type SceneCollection struct {
	*mgz.Core[*voiceagent.Scene]
}

func NewSceneCollection(db *mongo.Database) *SceneCollection {
	return &SceneCollection{
		Core: mgz.NewCore[*voiceagent.Scene](db, "va_scenes"),
	}
}
