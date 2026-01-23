package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"store/app/bi/internal/data/types"
	"store/pkg/clients/mgz"
)

type JuGuangCollection struct {
	*mgz.CoreCollection[types.JuGuangEvent]
}

func NewJuGuangCollection(db *mongo.Database) *JuGuangCollection {
	return &JuGuangCollection{
		CoreCollection: mgz.NewCoreCollection[types.JuGuangEvent](db, "juguang_events"),
	}
}
